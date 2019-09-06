package wechat_qrcode

import (
	"github.com/gin-gonic/gin"
	"dcx.com/tools/string_tool"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/xmdas-link/auth"
	"log"
	"time"
)

const Name = "wechat_qrcode"

type OAuthConfig struct {
	ClientID    string // APP的ClientID
	Secret      string // APP的密钥
	AuthBaseUrl string // 认证请求地址
	CallbackUrl string // APP授权的回调地址
	Scopes      []string
}

type Provider struct {
	init bool
	Name string
	DB   *gorm.DB
	*Client
	Session auth.Session
}

// 注册Provider时执行
func (p *Provider) OnProviderRegister(a *auth.GinAuth) error {

	if p.init {
		return nil
	}

	if p.ClientID == "" || p.Secret == "" {
		return errors.New("微信APP信息未配置")
	}

	if p.DB == nil {
		p.DB = a.Config.Core.DB
	}

	if p.DB == nil {
		return errors.New("未能从GinAuth获取数据库配置")
	} else {
		p.DB.AutoMigrate(&WechatUser{})
	}

	if p.Session == nil {
		p.Session = a.Config.Core.Session
	}

	if p.Session == nil {
		return errors.New("未能从GinAuth获取Session配置")
	}

	return nil
}

// 认证名称
func (p *Provider) GetName() string {
	return p.Name
}

// 登录引导
func (p *Provider) OnGuideLogin(c *gin.Context) error {

	var (
		state    = string_tool.GetRandomString(6)
		dataType = c.Query("type")
	)

	p.Session.StartSession(c, true)
	p.Session.SetValue(c, "wechat_state", state)

	if dataType == "json" {
		// 微信嵌入页面用
		c.Set("data", p.AuthCodeData(state))
		c.Set(auth.CtxKeyResultType, auth.RspTypeJSON)
	} else {
		// 默认直接跳转
		c.Set(auth.CtxKeyResultType, auth.RspTypeRedirect)
		c.Set("redirect", p.AuthCodeURL(state))
	}
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

	var (
		code  = c.Query("code")
		state = c.Query("state")
	)

	p.Session.StartSession(c, true)
	wechatState := p.Session.GetValueString(c, "wechat_state")
	if state == "" || wechatState != state {
		err = errors.New("state不匹配")
		return
	}

	if code == "" {
		err = errors.New("授权登录失败")
		return
	}

	token, errToken := p.GetAccessToken(code)
	if errToken != nil {
		log.Printf("[wechat_qrcode OnLoginCallback]%v", errToken)
		err = errors.New("获取Access Token失败")
		return
	}

	user, dbErr := p.GetUser(token)
	if dbErr != nil {
		log.Printf("[wechat_qrcode OnLoginCallback]%v", dbErr)
		err = errors.New("数据库获取用户信息失败")
		return
	}

	if !user.IsActive() {
		return nil, errors.New("账号状态不可登录")
	}

	wechatUserInfo, err := p.GetWechatUser(token)
	if err != nil {
		log.Printf("[wechat_qrcode OnLoginCallback]%v", err)
		err = errors.New("获取微信用户信息失败")
		return
	}

	expiredTime := time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
	user.TokenExpired = &expiredTime
	user.Name = wechatUserInfo.Nickname
	user.Token = token.AccessToken

	if dbErr := p.UpdateUser(user); dbErr != nil {
		log.Printf("[wechat_qrcode OnLoginCallback]%v", dbErr)
		return nil, errors.New("数据库更新户信息失败")
	}

	u = &UserData{
		WechatUser: &user,
		Provider:   p.GetName(),
	}
	return u, nil
}

func (p *Provider) GetUser(token *AccessToken) (data WechatUser, err error) {

	var (
		tx = p.DB
	)

	if token.Unionid != "" {
		err = tx.First(&data, "union_id = ?", token.Unionid).Error
	} else if token.Openid != "" {
		err = tx.First(&data, "open_id = ?", token.Openid).Error
	} else {
		err = errors.New("Token信息错误！")
	}

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			data.OpenId = token.Openid
			data.UnionId = token.Unionid
			data.Role = "common"
			err = tx.Create(&data).Error
		}
	}

	return
}

func (p *Provider) UpdateUser(data WechatUser) error {
	var (
		tx = p.DB
	)

	return tx.Model(&WechatUser{}).Where("id = ?", data.ID).Update(&WechatUser{
		Name:         data.Name,
		Token:        data.Token,
		TokenExpired: data.TokenExpired,
	}).Error
}
