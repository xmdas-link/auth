package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"net/http"
)

type ResultType string

const RspTypeJSON ResultType = "json"
const RspTypeString ResultType = "html_string"
const RspTypeTmpl ResultType = "tmpl"
const RspTypeRedirect ResultType = "redirect"

type Result struct {
	StatusCode int
	RspType    ResultType
	Value      interface{}
	Key        string
}

func (r *Result) GetStringValue() string {
	return r.Value.(string)
}

func (r *Result) Response(c *gin.Context) error {

	switch r.RspType {
	case RspTypeJSON:
		c.Render(r.StatusCode, render.JSON{Data: r.Value})
	case RspTypeRedirect:
		c.Redirect(r.StatusCode, r.GetStringValue())
	case RspTypeTmpl:
		c.HTML(r.StatusCode, r.Key, r.Value)
	case RspTypeString:
		c.Render(r.StatusCode, render.Data{
			ContentType: "text/html; charset=utf-8",
			Data:        []byte(r.GetStringValue()),
		})
	default:
		return fmt.Errorf("返回类型%v未定义", r.RspType)
	}

	return nil
}

func NewJSONResult(v interface{}) *Result {

	return &Result{
		StatusCode: http.StatusOK,
		RspType:    RspTypeJSON,
		Value:      v,
	}

}

func NewRedirectResult(url string) *Result {
	return &Result{
		StatusCode: http.StatusMovedPermanently,
		RspType:    RspTypeRedirect,
		Value:      url,
	}
}

func NewTmplResult(tplName string, v interface{}) *Result {
	return &Result{
		StatusCode: http.StatusOK,
		RspType:    RspTypeTmpl,
		Value:      v,
		Key:        tplName,
	}
}

func NewStringResult(v string) *Result {
	return &Result{
		StatusCode: http.StatusOK,
		RspType:    RspTypeString,
		Value:      v,
	}
}
