package auth

import (
	"github.com/gin-gonic/gin"
)

// 从Context里获取GinAuth
func GetAuthFromContext(c *gin.Context) *GinAuth {

	if obj, exist := c.Get(CtxKeyGinAuth); exist {
		if auth, ok := obj.(*GinAuth); ok {
			return auth
		}
	}

	return nil
}

//  从Context里获取默认的返回格式
func GetDefaultResultType(c *gin.Context) ResultType {

	if obj, exist := c.Get(CtxKeyResultType); exist {
		if resType, ok := obj.(ResultType); ok {
			return resType
		}
	}

	return ""
}

// 从Context里获取角色信息
func GetRoleFromContext(c *gin.Context) string {
	return c.GetString(CtxKeyUserRole)
}
