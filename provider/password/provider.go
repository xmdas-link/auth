package password

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mojocn/base64Captcha"
	"github.com/xmdas-link/auth"
	"golang.org/x/crypto/bcrypt"
)

const Name = "password"

type Provider struct {
	DB        *gorm.DB
	CryptCost int
	Name      string
	Captcha   *base64Captcha.ConfigCharacter
}

// 注册Provider时执行
func (p *Provider) OnProviderRegister(a *auth.GinAuth) error {

	if p.DB == nil && a.Config.Core.DB != nil {
		p.DB = a.Config.Core.DB
	}

	if p.DB != nil {
		p.DB.AutoMigrate(&User{})
	} else {
		return errors.New("需要从GinAuth获取数据库配置")
	}

	return nil
}

// 认证名称
func (p *Provider) GetName() string {
	return p.Name
}

// 登录引导
func (p *Provider) OnGuideLogin(c *gin.Context) error {
	if p.Captcha != nil {
		c.Set("captcha", true)
		// 需要使用验证码
		if c.Query("refresh") == "captcha" {
			idKeyC, capC := base64Captcha.GenerateCaptcha("", *p.Captcha)
			base64Png := base64Captcha.CaptchaWriteToBase64Encoding(capC)
			c.Set("data", map[string]interface{}{
				"captchaImg": base64Png,
				"captchaId":  idKeyC,
			})
			c.Set("refresh", "captcha")
		}
	}
	return nil
}

// 登录账号
func (p *Provider) OnLogin(c *gin.Context) (u auth.User, err error) {

	var (
		loginName = c.PostForm("user")
		pass      = c.PostForm("pass")
	)

	if p.Captcha != nil {
		err = p.VerifyCaptcha(c.PostForm("captcha_id"), c.PostForm("captcha"))
		if err != nil {
			return nil, err
		}
	}

	if loginName == "" {
		return nil, errors.New("请输入账号")
	}

	if pass == "" {
		return nil, errors.New("请输入密码")
	}

	user, dbErr := p.GetUser(loginName)
	if dbErr != nil {
		return nil, errors.New("账号不存在")
	}

	if err := p.ComparePassword(user.Password, pass); err != nil {
		return nil, errors.New("密码错误")
	}

	if !user.IsActive() {
		return nil, errors.New("账号状态不可登录")
	}

	u = &UserData{
		User:     user,
		Provider: p.GetName(),
	}

	return
}

// 登出账号
func (p *Provider) OnLogout(c *gin.Context) (u auth.User, err error) {
	return
}

// 第三方登录回调
func (p *Provider) OnLoginCallback(c *gin.Context) (u auth.User, err error) {
	return nil, fmt.Errorf("%s不支持第三方登录回调", p.GetName())
}

func (p *Provider) AddUser(user, pass, role, name string) (*User, error) {
	var (
		tx   = p.DB
		data = User{
			Name:      name,
			LoginName: user,
			Role:      role,
		}
		err error
	)

	if data.Password, err = p.EncryptPassword(pass); err != nil {
		return nil, err
	}

	if err := tx.Create(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func (p *Provider) GetUser(user string) (*User, error) {
	var (
		tx   = p.DB
		data = User{}
	)

	if err := tx.First(&data, "login_name = ?", user).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func (p *Provider) EncryptPassword(pass string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), p.CryptCost)
	return string(hashedPassword), err
}

func (p *Provider) ComparePassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (p *Provider) VerifyCaptcha(id, value string) error {
	if id == "" {
		return errors.New("缺少随机码ID")
	}
	if value == "" {
		return errors.New("需要验证随机码")
	}
	if base64Captcha.VerifyCaptcha(id, value) {
		return nil
	}
	return errors.New("随机码错误")
}
