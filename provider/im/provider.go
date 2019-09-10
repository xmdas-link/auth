package im

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/xmdas-link/auth"
	"github.com/xmdas-link/auth/user_store"
	"github.com/xmdas-link/tools/string_tool"
	"log"
	"time"
)

const Name = "im"

type OAuthConfig struct {
	ClientID      string // IM生成的APP的ClientID
	Secret        string // IM定义的APP的密钥
	MattermostUrl string // IM地址
	ApiVersion    string // IM的api版本
	CallbackUrl   string // APP授权的回调地址
	Scopes        []string
}

type Provider struct {
	DB *gorm.DB
	//BaseUrl    string
	//UrlVersion string
	Name      string
	Security  bool
	Session   auth.Session
	UserStore user_store.UserStoreInterface
	//oauthCfg   *oauth2.Config
	*OAuthConfig
	*Client
}

// 注册Provider时执行
func (p *Provider) OnProviderRegister(a *auth.GinAuth) error {
	/*if p.ClientID == "" || p.Secret == "" {
		return errors.New("缺少OAuth相关配置")
	}*/
	if p.DB == nil {
		p.DB = a.Config.Core.DB
	}

	if p.DB == nil {
		return errors.New("需要从GinAuth获取数据库配置")
	} else {
		p.DB.AutoMigrate(&ImUser{})
	}

	if p.Session == nil {
		p.Session = a.Config.Core.Session
	}

	if p.Security && p.Session == nil {
		return errors.New("未能从GinAuth获取Session配置")
	}

	if p.UserStore == nil {
		p.UserStore = a.Config.Core.UserStore
	}

	if p.UserStore == nil {
		return errors.New("未能从GinAuth获取UserStore配置")
	}

	return nil
}

// 认证名称
func (p *Provider) GetName() string {
	return p.Name
}

func (p *Provider) OnGuideLogin(c *gin.Context) error {
	var (
		state = string_tool.GetRandomString(6)
	)

	if p.Security {
		p.Session.StartSession(c, true)
		p.Session.SetValue(c, "im_state", state)
	}

	// log.Printf(p.AuthCodeURL(state))
	c.Set("redirect", p.AuthCodeURL(state))
	return nil
}

// 登录账号
func (p *Provider) OnLogin(c *gin.Context) (u auth.User, err error) {
	return nil, errors.New("不支持OnLogin接口")
}

// 登出账号
func (p *Provider) OnLogout(c *gin.Context) (u auth.User, err error) {
	return
}

// 第三方登录回调
func (p *Provider) OnLoginCallback(c *gin.Context) (u auth.User, err error) {

	if errMsg := c.Query("error"); errMsg != "" {
		return nil, fmt.Errorf("登录失败：%s", p.GetErrorReason(errMsg))
	}

	var (
		code  = c.Query("code")
		state = c.Query("state")
	)

	if p.Security {
		p.Session.StartSession(c, true)
		imState := p.Session.GetValueString(c, "im_state")
		if state == "" || imState != state {
			err = errors.New("state不匹配")
			return
		}
	}

	if code == "" {
		err = errors.New("授权登录失败")
		return
	}

	token, errToken := p.GetAccessToken(code)
	if errToken != nil {
		log.Printf("[im OnLoginCallback]%v", errToken)
		err = errors.New("获取Access Token失败")
		return
	}

	meData, imErr := p.GetMe(token)
	if imErr != nil {
		return nil, imErr
	}

	// 查找imUser
	imUser, dbErr := p.GetImUser(meData.ID)
	if dbErr != nil && gorm.IsRecordNotFoundError(dbErr) {
		// 需要创建新用户
		if imUser, dbErr = p.NewUser(meData, c); dbErr != nil {
			log.Printf("[im OnLoginCallback]%v", dbErr)
			err = errors.New("创建新用户失败！")
			return
		}
	}

	if imUser == nil || imUser.ID == 0 {
		return nil, errors.New("获取IM用户信息失败")
	}

	if imUser.UserBase, err = p.UserStore.Get(imUser.UID, c); err != nil {
		log.Printf("[im OnLoginCallback]%v", err)
		err = errors.New("获取UserStore内用户信息失败")
		return nil, err
	}

	if !imUser.UserBase.IsActive() {
		return nil, errors.New("账号状态不可登录")
	}

	// 更新用户信息
	imUser.Token = token.AccessToken

	// token到期时间
	if token.Expires > 0 {
		timeExpire := time.Unix(token.Expires, 0)
		imUser.TokenExpired = &timeExpire
	} else if token.ExpiresIn > 0 {
		timeExpire := time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
		imUser.TokenExpired = &timeExpire
	}

	if meData.Nickname != "" {
		imUser.Name = meData.Nickname
	} else {
		imUser.Name = meData.LastName + meData.FirstName
	}

	if dbErr := p.UpdateUser(imUser); dbErr != nil {
		log.Printf("[im OnLoginCallback]%v", dbErr)
		return nil, errors.New("数据库更新户信息失败")
	}

	// log.Printf("[im OnLoginCallback]头像地址：", p.GetHeaderUrl(imUser.ImID))

	u = &UserData{
		ImUser:    imUser,
		Provider:  p.GetName(),
		HeaderUrl: p.GetHeaderUrl(imUser.ImID),
	}
	return u, nil

}

func (p *Provider) GetImUser(id string) (*ImUser, error) {

	var (
		tx   = p.DB
		data = ImUser{}
	)

	err := tx.Find(&data, "im_id = ?", id).Error
	return &data, err
}

func (p *Provider) NewUser(imData ImUserData, c *gin.Context) (*ImUser, error) {

	var (
		newUserData = map[string]string{
			"name": imData.Nickname,
			"role": "common",
		}
		data = ImUser{
			ImID:      imData.ID,
			LoginName: imData.Username,
		}
		err error
	)

	if imData.Nickname != "" {
		newUserData["name"] = imData.Nickname
	} else {
		newUserData["name"] = imData.LastName + imData.FirstName
	}

	data.UserBase, data.UID, err = p.UserStore.New(newUserData, c)
	if err != nil {
		return nil, err
	} else if data.UID == "" {
		return nil, errors.New("获得新用户ID失败！")
	}

	if dbErr := p.DB.Create(&data).Error; dbErr != nil {
		return nil, dbErr
	}
	return &data, nil
}

func (p *Provider) UpdateUser(imUser *ImUser) error {
	var (
		tx = p.DB
	)

	return tx.Model(&ImUser{}).Where("id = ?", imUser.ID).Update(&ImUser{
		Name:         imUser.Name,
		Token:        imUser.Token,
		TokenExpired: imUser.TokenExpired,
	}).Error
}

func (Provider) GetErrorReason(v string) string {
	switch v {
	case "access_denied":
		return "用户拒绝了登录授权"
	default:
		return v
	}
}
