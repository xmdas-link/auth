package auth

import (
	"github.com/gin-gonic/gin"
)

func DefaultNotFoundHandler(c *gin.Context) {

	var (
		resType = GetDefaultResultType(c)
		result  *Result
	)

	switch resType {
	case RspTypeTmpl:
		result = NewTmplResult("provider_not_found.tmpl", gin.H{})
	case RspTypeString:
		result = NewStringResult("<!DOCTYPE html><html lang=\"zh-CN\"><head><meta charset=\"UTF-8\"><title>404错误</title></head><body><p>登录方式不存在</p></body></html>")
	case RspTypeRedirect:
		a := GetAuthFromContext(c)
		url := "/"
		if a != nil && a.Config.Path.RootURL != "" {
			url = a.Config.Path.RootURL
		}
		result = NewRedirectResult(url)
	case RspTypeJSON:
		// 默认返回值
		fallthrough
	default:
		result = NewJSONResult(gin.H{"code": 404, "message": "登录方式不存在！"})
	}

	result.Response(c)

}
