package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/xmdas-link/auth"
	"github.com/xmdas-link/auth/provider/im"
	"github.com/xmdas-link/auth/provider/password"
	"github.com/xmdas-link/auth/provider/wechat_qrcode"
	"github.com/xmdas-link/auth/redis_token"
	"github.com/xmdas-link/auth/render/oauth_render"
	"github.com/xmdas-link/auth/render/password_render"
	"github.com/xmdas-link/auth/session"
	"net/http"
	"path/filepath"
	"time"
)

func main() {

	var (
		gAuth             *auth.GinAuth
		redisToken        *redis_token.Module
		cfg               = auth.Config{}
		router            = gin.Default()
		db, err           = gorm.Open("mysql", "root:@(127.0.0.1:3306)/gin_test?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai")
		publicKey, _      = filepath.Abs("example/keys/token_key.pub")
		privateKye, _     = filepath.Abs("example/keys/token_key")
		templateFolder, _ = filepath.Abs("example/templates/*")
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

	// 配置token（使用redis存储token）
	redisCfg := &redis_token.Config{}
	redisCfg.Redis.Address = "192.168.0.70:6379"
	redisCfg.Redis.Db = 15
	redisCfg.Jwt.ExpireDuration = time.Hour * 24 // Token 默认有效期
	redisCfg.Jwt.PublicKeyPath = publicKey
	redisCfg.Jwt.PrivateKeyPath = privateKye

	if redisToken, err = redis_token.New(redisCfg); err != nil {
		panic(err)
	}

	// Session配置
	sessService, sessErr := session.New(session.Config{
		RedisAddress: "192.168.0.70:6379",
		AuthKey:      "111",
	})
	if sessErr != nil {
		panic(sessErr)
	}

	// 配置登录
	cfg.Path.Mount = "/auth"
	cfg.Path.TokenKey = "example_jwt"
	cfg.Path.Domain = "127.0.0.1"
	cfg.Path.RedirectAfterLogin = "http://127.0.0.1:9528"
	cfg.Core.Router = router
	cfg.Core.DB = db
	cfg.Core.AuthToken = redisToken
	cfg.Core.Session = sessService

	gAuth, err = auth.New(cfg)
	if err != nil {
		panic(err)
	}

	// 注册登录方式
	passProvider := password.New()

	// 账号登录渲染
	passRender := &password_render.Render{}
	if regErr := gAuth.RegisterProvider(passProvider, passRender); regErr != nil {
		panic(regErr)
	}

	// IM登录
	imProvider := im.New(&im.OAuthConfig{
		ClientID:      "your_client",
		Secret:        "your_wecret",
		MattermostUrl: "http://your.domain.com",
		CallbackUrl:   "http://your.domain.com:8009/auth/login_callback/im",
	})

	// oauth2.0渲染
	oauthRender := &oauth_render.Render{}
	if regErr := gAuth.RegisterProvider(imProvider, oauthRender); regErr != nil {
		panic(regErr)
	}

	// 微信扫码登录
	wxProvider := wechat_qrcode.New(&wechat_qrcode.OAuthConfig{
		ClientID:    "your_client",
		Secret:      "your_wecret",
		CallbackUrl: "http://your.domain.com:8009/auth/login_callback/wechat_qrcode",
		Scopes:      []string{"snsapi_login"},
	})
	if regErr := gAuth.RegisterProvider(wxProvider, oauthRender); regErr != nil {
		panic(regErr)
	}

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
