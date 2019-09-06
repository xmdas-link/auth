package password_render

import (
	"github.com/gin-gonic/gin"
	"github.com/xmdas-link/auth"
)

type Render struct {
	init bool
}

func (r *Render) OnRenderRegister(a *auth.GinAuth) error {
	if r.init {
		return nil
	}
	r.init = true
	return nil
}

func (r *Render) Error(c *gin.Context) (ret *auth.Result, err error) {
	var (
		code = c.GetInt("code")
		msg  = c.GetString("err")
	)

	if msg == "" {
		msg = "错误信息未定义"
	}

	ret = auth.NewJSONResult(gin.H{"code": code, "message": msg})
	return
}

func (r *Render) GuideLogin(c *gin.Context) (ret *auth.Result, err error) {
	ret = auth.NewTmplResult("password_login.tmpl", gin.H{})
	return
}

func (r *Render) FailLogin(c *gin.Context) (ret *auth.Result, err error) {
	ret, err = r.Error(c)
	return
}

func (r *Render) SuccessLogin(c *gin.Context, u auth.User) (ret *auth.Result, err error) {

	var (
		a = auth.GetAuthFromContext(c)
	)

	if a == nil {
		c.Set("err", "GinAuth missing！")
		return r.Error(c)
	}

	tokenInfo := map[string]interface{}{
		"key":      a.Config.Path.TokenKey,
		"token":    u.GetToken(),
		"expired":  u.GetExpired(),
		"domain":   c.GetString("domain"),
		"redirect": c.GetString("redirect"),
	}

	ret = auth.NewJSONResult(gin.H{"code": 1, "data": tokenInfo})
	return
}

func (r *Render) Logout(c *gin.Context) (ret *auth.Result, err error) {
	if _, hasErr := c.Get("err"); hasErr {
		return r.Error(c)
	}

	if redirect := c.GetString("redirect"); redirect == "" {
		ret = auth.NewJSONResult(gin.H{"code": 1, "message": "已退出"})
	} else {
		ret = auth.NewRedirectResult(redirect)
	}
	return
}
