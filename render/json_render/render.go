package json_render

import (
	"github.com/gin-gonic/gin"
	"github.com/xmdas-link/auth"
)

type Render struct {
	init    bool
	LogFunc func(c *gin.Context, Action string, result interface{})
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

	result := gin.H{"code": code, "message": msg}

	ret = auth.NewJSONResult(result)
	if r.LogFunc != nil {
		r.LogFunc(c, "Error", result)
	}
	return
}

func (r *Render) GuideLogin(c *gin.Context) (ret *auth.Result, err error) {

	// 检查是不是有错误信息
	if _, hasErr := c.Get("err"); hasErr {
		return r.Error(c)
	}

	data, _ := c.Get("data")
	result := gin.H{"code": 1, "data": data}

	ret = auth.NewJSONResult(result)
	if r.LogFunc != nil {
		r.LogFunc(c, "GuideLogin", result)
	}
	return
}

func (r *Render) FailLogin(c *gin.Context) (ret *auth.Result, err error) {
	return r.Error(c)
}

func (r *Render) SuccessLogin(c *gin.Context, u auth.User) (ret *auth.Result, err error) {

	data := map[string]interface{}{
		"token":   u.GetToken(),
		"expired": u.GetExpired(),
	}

	result := gin.H{"code": 1, "data": data, "message": ""}
	ret = auth.NewJSONResult(result)
	if r.LogFunc != nil {
		r.LogFunc(c, "SuccessLogin", result)
	}
	return
}

func (r *Render) Logout(c *gin.Context) (ret *auth.Result, err error) {
	if _, hasErr := c.Get("err"); hasErr {
		return r.Error(c)
	}
	result := gin.H{"code": 1, "message": "已退出"}
	ret = auth.NewJSONResult(result)
	if r.LogFunc != nil {
		r.LogFunc(c, "Logout", result)
	}
	return
}
