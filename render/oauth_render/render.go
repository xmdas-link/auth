package oauth_render

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xmdas-link/auth"
)

type Render struct {
	init bool
	// 有配置redirectMap值的情况，登录后跳转优先用匹配的跳转地址，没有才会使用默认的跳转定义
	RedirectMap map[string]string
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
	var (
		resultType, _ = c.Get(auth.CtxKeyResultType)
	)

	if resultType == auth.RspTypeJSON {
		if data, exist := c.Get("data"); exist {
			ret = auth.NewJSONResult(gin.H{"code": 1, "data": data})
			return
		}
	}

	url := c.GetString("redirect")
	if url == "" {
		return nil, errors.New("Provider未设置登录跳转地址redirect")
	}
	ret = auth.NewRedirectResult(url)
	return
}

func (r *Render) FailLogin(c *gin.Context) (ret *auth.Result, err error) {
	if redirect := c.GetString("redirect"); redirect != "" {
		return auth.NewRedirectResult(redirect), nil
	}
	ret, err = r.Error(c)
	return
}

func (r *Render) SuccessLogin(c *gin.Context, u auth.User) (ret *auth.Result, err error) {

	var (
		a        = auth.GetAuthFromContext(c)
		redirect = c.GetString("redirect")
	)

	if a == nil {
		c.Set("err", "GinAuth missing！")
		return r.Error(c)
	}

	if key := c.GetString("from"); key != "" && r.RedirectMap != nil {
		if overrideUrl, exist := r.RedirectMap[key]; exist {
			redirect = overrideUrl
		}
	}

	ret = auth.NewTmplResult("oauth_login.tmpl", gin.H{
		"key":      a.Config.Path.TokenKey,
		"token":    u.GetToken(),
		"expired":  u.GetExpired(),
		"domain":   c.GetString("domain"),
		"redirect": redirect,
	})
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
