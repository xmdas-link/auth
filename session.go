package auth

import "github.com/gin-gonic/gin"

type Session interface {
	// 对一个request启动session，默认要检查csrf
	StartSession(c *gin.Context, noCheckCsrf bool)
	// 从session获取值
	GetValue(c *gin.Context, key string) (v interface{}, exist bool)
	// 按字符串格式读取一个session值
	GetValueString(c *gin.Context, key string) string
	// 写入session
	SetValue(c *gin.Context, key string, v interface{})
	// 移除值
	DelValue(c *gin.Context, key string)
	// 清空session
	ClearSession(c *gin.Context)
}
