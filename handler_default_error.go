package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
)

var (
	ErrMissingResult = errors.New("接口未按要求返回结果")
)

func DefaultErrorResponse(c *gin.Context, err error) {
	var (
		resType = GetDefaultResultType(c)
		result  *Result
	)

	switch resType {
	case RspTypeTmpl:
		result = NewTmplResult("error.tmpl", gin.H{"error": err.Error()})
	case RspTypeString:
		result = NewStringResult("<!DOCTYPE html><html lang=\"zh-CN\"><head><meta charset=\"UTF-8\"><title>500错误</title></head><body><p>执行发生错误：" + err.Error() + "</p></body></html>")
	default:
		DefaultJsonError(c, err.Error())
		return
	}

	result.Response(c)

}

func DefaultJsonError(c *gin.Context, errMsg string) {
	result := NewJSONResult(gin.H{"code": 0, "message": errMsg})
	result.Response(c)
}

func DefaultJsonSuccess(c *gin.Context, data interface{}) {
	result := NewJSONResult(gin.H{"code": 1, "message": "", "data": data})
	result.Response(c)
}
