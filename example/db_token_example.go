package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/mojocn/base64Captcha"
	"github.com/xmdas-link/auth"
	"github.com/xmdas-link/auth/db_token"
	"github.com/xmdas-link/auth/provider/im"
	"github.com/xmdas-link/auth/provider/password"
	"github.com/xmdas-link/auth/render/oauth_render"
	"github.com/xmdas-link/auth/render/password_render"
	"net/http"
	"path/filepath"
	"time"
)

func main() {

	var (
		gAuth             *auth.GinAuth
		dbToken           *db_token.Module
		cfg               = auth.Config{}
		router            = gin.Default()
		db, err           = gorm.Open("mysql", "root:@(127.0.0.1:3306)/gin_test?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai")
		publicKey, _      = filepath.Abs("keys/rsa_public_key.pem")
		privateKey, _     = filepath.Abs("keys/rsa_private_key.pem")
		templateFolder, _ = filepath.Abs("templates/*")
	)

	if err != nil {
		panic(err)
	}

	// 打开sql记录
	db.LogMode(true)
	//router.Use(middleware.NoCache)
	//router.Use(middleware.Options)
	//router.Use(middleware.Cross)

	// 定义router的模板文件地址
	router.LoadHTMLGlob(templateFolder)

	// 配置token（使用数据库存储token）
	if dbToken, err = db_token.New(&db_token.Config{
		DB:             db,
		ExpireDuration: time.Hour * 24, // Token 默认有效期
		PublicKeyPath:  publicKey,
		PrivateKeyPath: privateKey,
	}); err != nil {
		panic(err)
	}

	// 配置登录
	cfg.Path.Mount = "/auth"
	cfg.Path.TokenKey = "example_jwt"
	cfg.Path.Domain = "127.0.0.1"
	cfg.Path.RedirectAfterLogin = "http://127.0.0.1:9528"
	cfg.Core.Router = router
	cfg.Core.DB = db
	cfg.Core.AuthToken = dbToken

	gAuth, err = auth.New(cfg)
	if err != nil {
		panic(err)
	}

	// 注册登录方式
	passProvider := password.New()
	// 配置验证码
	passProvider.Captcha = &base64Captcha.ConfigCharacter{
		Height: 40,
		Width:  200,
		//const CaptchaModeNumber:数字,CaptchaModeAlphabet:字母,CaptchaModeArithmetic:算术,CaptchaModeNumberAlphabet:数字字母混合.
		Mode:               base64Captcha.CaptchaModeArithmetic,
		ComplexOfNoiseText: base64Captcha.CaptchaComplexLower,
		ComplexOfNoiseDot:  base64Captcha.CaptchaComplexLower,
		IsUseSimpleFont:    true,
		IsShowHollowLine:   false,
		IsShowNoiseDot:     false,
		IsShowNoiseText:    true,
		IsShowSlimeLine:    false,
		IsShowSineLine:     false,
		CaptchaLen:         4,
	}

	// 账号登录渲染
	passRender := &password_render.Render{}
	if regErr := gAuth.RegisterProvider(passProvider, passRender); regErr != nil {
		panic(regErr)
	}

	// IM登录
	imProvider := im.New(&im.OAuthConfig{
		ClientID:      "your_client_id",
		Secret:        "your_secret",
		MattermostUrl: "http://your.domain",
		CallbackUrl:   "http://192.168.0.88:8009/auth/login_callback/im",
	})

	// oauth2.0渲染
	oauthRender := &oauth_render.Render{}
	if regErr := gAuth.RegisterProvider(imProvider, oauthRender); regErr != nil {
		panic(regErr)
	}

	// 微信扫码登录
	/*wxProvider := wechat_qrcode.New(&wechat_qrcode.OAuthConfig{
		ClientID:    "your_client_id",
		Secret:      "your_secret",
		CallbackUrl: "http://your.domain/auth/login_callback/wechat_qrcode",
		Scopes:      []string{"snsapi_login"},
	})
	if regErr := gAuth.RegisterProvider(wxProvider, oauthRender); regErr != nil {
		panic(regErr)
	}*/

	if mountErr := gAuth.MountAuth(); mountErr != nil {
		panic(mountErr)
	}

	// 镶入式微信扫码登录的例子
	router.GET("/wechat_qrcode", func(c *gin.Context) {
		c.HTML(http.StatusOK, "wechat_qrcode_login.tmpl", gin.H{})
	})

	////////////////////////////////需要认证身份分割线///////////////////////////////////

	// 挂认证的中间件，挂中间件后的所有连接都要走认证
	router.Use(gAuth.AuthTokenMiddleware())

	router.GET("/user/rights", func(c *gin.Context) {

		user := c.GetStringMapString(auth.CtxKeyAuthUser)
		userRole := auth.GetRoleFromContext(c)

		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "I am a test", "data": gin.H{"role": userRole, "user": user, "rights": []string{"/user/rights"}}})
	})

	router.Run(":8009")

}
