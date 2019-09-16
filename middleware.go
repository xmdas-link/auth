package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
)

func (a *GinAuth) AuthTokenMiddleware() func(c *gin.Context) {

	return func(c *gin.Context) {

		var (
			token = c.GetHeader(a.Config.Path.HeaderKey)
		)

		if token == "" {
			DefaultErrorResponse(c, errors.New("Token is required!"))
			c.Abort()
			return
		}

		authToken := a.Config.Core.AuthToken
		user := authToken.FindToken(token)
		if user == nil {
			c.Set(CtxKeyErrorCode, DefaultErrorTokenInvalid)
			DefaultErrorResponse(c, errors.New("Token invalid!"))
			c.Abort()
			return
		}

		if uid, exist := user["uid"]; exist && a.Config.Core.UserStore != nil {
			userBase, err := a.Config.Core.UserStore.Get(uid, c)
			if err != nil {
				c.Set(CtxKeyErrorCode, DefaultErrorTokenInvalid)
				DefaultErrorResponse(c, errors.New("User not found!"))
				c.Abort()
				return
			}
			if !userBase.IsActive() {
				c.Set(CtxKeyErrorCode, DefaultErrorTokenInvalid)
				DefaultErrorResponse(c, errors.New("Account not active!"))
				c.Abort()
				return
			}

			user["role"] = userBase.GetRole()
		}

		c.Set(CtxKeyUserRole, user["role"])
		c.Set(CtxKeyAuthUser, user)

		c.Next()
	}

}
