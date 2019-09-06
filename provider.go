package auth

import (
	"github.com/gin-gonic/gin"
)

// 认证必须实现的接口
type AuthProvider interface {

	// 认证名称
	GetName() string

	// 注册时需要执行的
	OnProviderRegister(a *GinAuth) error

	// 登录引导
	OnGuideLogin(c *gin.Context) error

	// 登录账号
	OnLogin(c *gin.Context) (User, error)
	// 第三方登录回调
	OnLoginCallback(c *gin.Context) (User, error)
	OnLogout(c *gin.Context) (User, error)
}

type User interface {
	GetID() string
	GetProvider() string
	GetRole() string
	GetToken() string
	GetExpired() int64
	SetToken(v string, expiredAt int64)
	GetMapData() map[string]string
}
