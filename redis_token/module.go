package redis_token

import (
	"crypto/md5"
	"crypto/rsa"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"io/ioutil"
	"log"
	"time"
)

type Config struct {
	Redis struct {
		Address  string
		Password string
		Db       int
	}
	Jwt struct {
		ExpireDuration time.Duration
		PrivateKeyPath string
		PublicKeyPath  string
		JwtSignMethod  string
	}
}

type Module struct {
	client    *redis.Client
	signKey   *rsa.PrivateKey
	verifyKey *rsa.PublicKey
	*Config
}

func New(cfg *Config) (*Module, error) {

	if cfg.Redis.Address == "" {
		return nil, errors.New("Redis配置未填写")
	}

	var (
		m = &Module{
			client: redis.NewClient(&redis.Options{
				Addr:     cfg.Redis.Address,
				Password: cfg.Redis.Password,
				DB:       cfg.Redis.Db,
			}),
			Config: cfg,
		}
	)

	if cfg.Jwt.JwtSignMethod == "" {
		cfg.Jwt.JwtSignMethod = "RS256"
	}

	if signBytes, err := ioutil.ReadFile(cfg.Jwt.PrivateKeyPath); err != nil {
		return nil, err
	} else if m.signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes); err != nil {
		return nil, err
	}

	if verifyBytes, err := ioutil.ReadFile(cfg.Jwt.PublicKeyPath); err != nil {
		return nil, err
	} else if m.verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes); err != nil {
		return nil, err
	}

	return m, nil
}

type UserClaims struct {
	*jwt.StandardClaims
	User map[string]string
}

// 新Token
func (m *Module) NewToken(user map[string]string) (string, int64, error) {

	var (
		t         = jwt.New(jwt.GetSigningMethod(m.Jwt.JwtSignMethod))
		expiredAt = time.Now().Add(m.Jwt.ExpireDuration).Unix()
		uClaims   = &UserClaims{
			&jwt.StandardClaims{
				ExpiresAt: expiredAt,
			},
			user,
		}
		token string
		err   error
	)

	t.Claims = uClaims
	token, err = t.SignedString(m.signKey)
	if err == nil {
		// 保存token
		err = m.StoreToken(uClaims, token)
	}
	if err != nil {
		log.Printf("ERROR: [redis_token.NewToken]%v", err)
	}

	return m.EncodeToken(token), expiredAt, err
}

// 清除Token
func (m *Module) ClearToken(tokenMd5 string) error {
	return m.client.Del(tokenMd5).Err()
}

// 清除用户的Token
func (m *Module) ClearTokenOfUser(uid string, provider string) error {
	var (
		key = m.GetUserTokenKey(uid, provider)
	)
	if oldToken, err := m.client.Get(key).Result(); err == nil && oldToken != "" {
		return m.ClearToken(oldToken)
	}
	return nil
}

// 查找token
func (m *Module) FindToken(tokenMd5 string) (user map[string]string) {

	jwtToken, err := m.client.Get(tokenMd5).Result()
	if err != nil {
		log.Printf("ERROR: [redis_token FindToken]%v", err)
		return nil
	}

	tokenObj, err := jwt.ParseWithClaims(jwtToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return m.verifyKey, nil
	})

	if err == nil {
		claims := tokenObj.Claims.(*UserClaims)
		if claims.ExpiresAt > time.Now().Unix() {
			return claims.User
		}
	} else {
		log.Printf("ERROR: [redis_token.FindToken]%v", err)
	}

	return nil
}

func (m *Module) StoreToken(claim *UserClaims, token string) error {
	var (
		user     = claim.User
		tokenMd5 = m.EncodeToken(token)
	)

	if user["id"] == "" || user["provider"] == "" || user["ip"] == "" {
		return errors.New("User缺少必要的字段")
	}

	log.Printf("StoreToken:len(%d):%v", len(token), token)

	// 清除旧token
	m.ClearTokenOfUser(user["id"], user["provider"])

	// 写新token
	m.client.Set(m.GetUserTokenKey(user["id"], user["provider"]), tokenMd5, m.Jwt.ExpireDuration)
	return m.client.Set(tokenMd5, token, m.Jwt.ExpireDuration).Err()
}

func (m *Module) EncodeToken(token string) string {
	h := md5.New()
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}

func (m *Module) GetUserTokenKey(uid string, provider string) string {
	return fmt.Sprintf("%v:%v", uid, provider)
}
