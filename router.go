package auth

import (
	"errors"
)

func (a *GinAuth) MountAuth() error {

	if a.Config.Core.Router == nil {
		return errors.New("Router未定义！")
	}

	var (
		g      = a.Config.Core.Router
		rGroup = g.Group(a.Config.Path.Mount)
	)

	// 模板加载方式需要重构，目测不会放在这里了
	/*if a.Config.Response.DefaultResultType == RspTypeTmpl && a.Config.Response.TemplateFolder == "" {
		return errors.New("默认返回方式为Tmpl渲染，需要定义TemplateFolder")
	}*/

	if a.Handlers.GetLoginHandler == nil {
		a.Handlers.GetLoginHandler = DefaultGetLoginHandler
	}
	rGroup.GET("/login/:name", NewGinHandler(a.Handlers.GetLoginHandler, a))

	if a.Handlers.PostLoginHandler == nil {
		a.Handlers.PostLoginHandler = DefaultPostLoginHandler
	}
	rGroup.POST("/login/:name", NewGinHandler(a.Handlers.PostLoginHandler, a))

	if a.Handlers.LogoutHandler == nil {
		a.Handlers.LogoutHandler = DefaultLogoutHandler
	}
	rGroup.GET("/logout", NewHandler(a.Handlers.LogoutHandler, a))
	rGroup.POST("/logout", NewHandler(a.Handlers.LogoutHandler, a))

	if a.Handlers.LoginCallbackHandler == nil {
		a.Handlers.LoginCallbackHandler = DefaultLoginCallbackHandler
	}
	rGroup.GET("/login_callback/:name", NewGinHandler(a.Handlers.LoginCallbackHandler, a))
	rGroup.POST("/login_callback/:name", NewGinHandler(a.Handlers.LoginCallbackHandler, a))

	return nil
}
