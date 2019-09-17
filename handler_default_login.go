package auth

import (
	"github.com/gin-gonic/gin"
)

func DefaultGetLoginHandler(c *gin.Context, p AuthProvider, r AuthRender) {

	if err := p.OnGuideLogin(c); err != nil {
		ret, err := r.Error(c)
		doResponse(c, ret, err)
		return
	}

	ret, err := r.GuideLogin(c)
	doResponse(c, ret, err)

}

func DefaultPostLoginHandler(c *gin.Context, p AuthProvider, r AuthRender) {

	var (
		user, err = p.OnLogin(c)
		hasErr    = false
	)

	if err != nil {
		c.Set("err", err.Error())
		hasErr = true
	} else if user == nil {
		c.Set("err", ErrMissingResult.Error())
		hasErr = true
	}

	if hasErr {
		ret, err := r.FailLogin(c)
		doResponse(c, ret, err)
	} else {
		AfterLogin(c, r, user)
	}

}

func DefaultLoginCallbackHandler(c *gin.Context, p AuthProvider, r AuthRender) {

	var (
		user, err = p.OnLoginCallback(c)
		hasErr    = false
	)

	if err != nil {
		c.Set("err", err.Error())
		hasErr = true
	} else if user == nil {
		c.Set("err", ErrMissingResult.Error())
		hasErr = true
	}

	if hasErr {
		ret, err := r.FailLogin(c)
		doResponse(c, ret, err)
	} else {
		AfterLogin(c, r, user)
	}
}

func AfterLogin(c *gin.Context, r AuthRender, user User) {
	var (
		gAuth   = GetAuthFromContext(c)
		token   string
		expired int64
		err     error
	)

	if gAuth == nil {
		c.Set("err", "GinAuth missing！")
		ret, err := r.FailLogin(c)
		doResponse(c, ret, err)
		return
	}

	userData := user.GetMapData()
	userData["id"] = user.GetID()
	userData["provider"] = user.GetProvider()
	userData["role"] = user.GetRole()
	userData["ip"] = c.ClientIP()

	token, expired, err = gAuth.Config.Core.AuthToken.NewToken(userData)
	if err != nil {
		c.Set("err", "Token生成失败！")
		ret, err := r.FailLogin(c)
		doResponse(c, ret, err)
		return
	}

	c.Set("domain", gAuth.Config.Path.Domain)
	c.Set("redirect", gAuth.Config.Path.RedirectAfterLogin)

	user.SetToken(token, expired)
	ret, err := r.SuccessLogin(c, user)
	doResponse(c, ret, err)
}

func doResponse(c *gin.Context, ret *Result, err error) {
	if err != nil {
		DefaultErrorResponse(c, err)
		return
	}

	// 如果ret为空，则认为已在render中处理了输出
	if ret != nil {
		ret.Response(c)
		//DefaultErrorResponse(c, ErrMissingResult)
		//return
	}

}
