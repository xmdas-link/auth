package auth

import "github.com/gin-gonic/gin"

type AuthRender interface {

	// 注册时需要执行的
	OnRenderRegister(a *GinAuth) error

	// 错误渲染
	Error(c *gin.Context) (*Result, error)

	// 登录引导
	GuideLogin(c *gin.Context) (*Result, error)

	// 错误登录
	FailLogin(c *gin.Context) (*Result, error)

	// 成功登录
	SuccessLogin(c *gin.Context, u User) (*Result, error)

	// 登出账号
	Logout(c *gin.Context) (*Result, error)
}
