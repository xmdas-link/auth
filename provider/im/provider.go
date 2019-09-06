package im

import (
	"context"
	"github.com/gin-gonic/gin"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/xmdas-link/auth"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"time"
)

const Name = "im"

type OAuthConfig struct {
	ClientID      string // IM生成的APP的ClientID
	Secret        string // IM定义的APP的密钥
	MattermostUrl string // IM地址
	CallbackUrl   string // APP授权的回调地址
}

type Provider struct {
	DB         *gorm.DB
	BaseUrl    string
	UrlVersion string
	Name       string
	oauthCfg   *oauth2.Config
	//*OAuthConfig
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

	return nil
}

// 认证名称
func (p *Provider) GetName() string {
	return p.Name
}

func (p *Provider) OnGuideLogin(c *gin.Context) error {
	url := p.oauthCfg.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Set("redirect", url)
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

	if code := c.Query("code"); code != "" {

		client, authErr := p.NewClient(code)
		if authErr != nil {
			return nil, authErr
		}

		meData, err := client.GetMe()
		if err != nil {
			return nil, err
		}

		imUser, dbErr := p.GetUser(meData)
		if dbErr != nil {
			log.Printf("[im OnLoginCallback]%v", dbErr)
			return nil, errors.New("数据库读取用户信息失败")
		}

		if !imUser.IsActive() {
			return nil, errors.New("账号状态不可登录")
		}

		// 更新用户信息
		imUser.Token, imUser.TokenExpired = client.GetAccessToken()

		if meData.Nickname != "" {
			imUser.Name = meData.Nickname
		} else {
			imUser.Name = meData.LastName + meData.FirstName
		}

		if dbErr := p.UpdateUser(imUser); dbErr != nil {
			log.Printf("[im OnLoginCallback]%v", dbErr)
			return nil, errors.New("数据库更新户信息失败")
		}

		u = &UserData{
			ImUser:   imUser,
			Provider: p.GetName(),
		}
		return u, nil
	}
	return nil, errors.New("缺少必需的code参数！")
}

func (p *Provider) NewClient(code string) (*Client, error) {

	httpClient := &http.Client{Timeout: 2 * time.Second}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)

	token, err := p.oauthCfg.Exchange(ctx, code)

	if err != nil {
		return nil, err
	}

	return &Client{
		BaseURL:    p.BaseUrl,
		URLVersion: p.UrlVersion,
		Token:      token,
		Im:         p.oauthCfg.Client(ctx, token),
	}, nil

}

func (p *Provider) GetUser(imData ImUserData) (*ImUser, error) {

	var (
		tx   = p.DB
		data = ImUser{}
		err  error
	)

	// 获取用户如果存在
	err = tx.Find(&data, "im_id = ?", imData.ID).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			data.ImID = imData.ID
			data.LoginName = imData.Username
			data.Role = "common"
			err = tx.Create(&data).Error
		}
	}

	if err != nil {
		return nil, err
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
