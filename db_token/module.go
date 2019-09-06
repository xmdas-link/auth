package db_token

import (
	"crypto/md5"
	"crypto/rsa"
	"encoding/hex"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"log"
	"time"
)

type Config struct {
	DB             *gorm.DB
	ExpireDuration time.Duration
	PrivateKeyPath string
	PublicKeyPath  string
	JwtSignMethod  string
}

type Module struct {
	*Config
	signKey   *rsa.PrivateKey
	verifyKey *rsa.PublicKey
}

func New(cfg *Config) (*Module, error) {

	var (
		m = &Module{
			Config: cfg,
		}
	)

	if m.JwtSignMethod == "" {
		m.JwtSignMethod = "RS256"
	}

	if signBytes, err := ioutil.ReadFile(m.PrivateKeyPath); err != nil {
		return nil, err
	} else if m.signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes); err != nil {
		return nil, err
	}

	if verifyBytes, err := ioutil.ReadFile(m.PublicKeyPath); err != nil {
		return nil, err
	} else if m.verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes); err != nil {
		return nil, err
	}

	if m.DB != nil {
		m.DB.AutoMigrate(&AuthUserToken{})
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
		t         = jwt.New(jwt.GetSigningMethod(m.JwtSignMethod))
		expiredAt = time.Now().Add(m.ExpireDuration).Unix()
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
		log.Printf("ERROR: [db_token.NewToken]%v", err)
	}

	return token, expiredAt, err
}

// 清除Token
func (m *Module) ClearToken(token string) error {
	var (
		tx = m.DB
	)

	return tx.Delete(&AuthUserToken{}, "token = ?", token).Error
}

// 清除用户的Token
func (m *Module) ClearTokenOfUser(uid string, provider string) error {
	var (
		tx = m.DB
	)

	return tx.Delete(&AuthUserToken{}, "uid = ? AND provider = ?", uid, provider).Error
}

// 查找token
func (m *Module) FindToken(token string) (user map[string]string) {

	var (
		tx       = m.DB
		count    = 0
		tokenMd5 = m.EncodeToken(token)
	)

	if dbErr := tx.Model(&AuthUserToken{}).Where("token = ? AND token_info = ?", tokenMd5, token).Count(&count).Error; dbErr != nil {
		log.Printf("ERROR: [db_token.FindToken]%v", dbErr)
		return nil
	} else if count == 0 {
		// token不存在
		return nil
	}

	tokenObj, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return m.verifyKey, nil
	})

	if err == nil {
		claims := tokenObj.Claims.(*UserClaims)
		if claims.ExpiresAt > time.Now().Unix() {
			return claims.User
		}
	} else {
		log.Printf("ERROR: [db_token.FindToken]%v", err)
	}

	return nil
}

func (m *Module) StoreToken(claim *UserClaims, token string) error {
	var (
		tx       = m.DB
		user     = claim.User
		tokenMd5 = m.EncodeToken(token)
	)

	if user["id"] == "" || user["provider"] == "" || user["ip"] == "" {
		return errors.New("User缺少必要的字段")
	}

	log.Printf("StoreToken:len(%d):%v", len(token), token)

	// 创建或覆盖
	if m.ExistUser(user["id"], user["provider"]) {
		return tx.Model(AuthUserToken{}).Where("uid = ? AND provider = ?", user["id"], user["provider"]).Update(map[string]interface{}{
			"ip":         user["ip"],
			"token":      tokenMd5,
			"token_info": token,
			"expired_at": claim.ExpiresAt,
		}).Error
	}

	return tx.Create(&AuthUserToken{
		Uid:       user["id"],
		Provider:  user["provider"],
		IP:        user["ip"],
		Token:     tokenMd5,
		TokenInfo: token,
		ExpiredAt: claim.ExpiresAt,
	}).Error

}

func (m *Module) ExistUser(id string, provider string) bool {
	var (
		tx    = m.DB
		count = 0
	)

	tx.Model(AuthUserToken{}).Where("uid = ? AND provider = ?", id, provider).Count(&count)

	return count > 0
}

func (m *Module) EncodeToken(token string) string {
	h := md5.New()
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}
