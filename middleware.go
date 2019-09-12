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
			err := errors.New("Token is required!")
			DefaultErrorResponse(c, err)
			c.Abort()
			return
		}

		authToken := a.Config.Core.AuthToken
		user := authToken.FindToken(token)
		if user == nil {
			err := errors.New("Token invalid!")
			DefaultErrorResponse(c, err)
			c.Abort()
			return
		}

		if uid, exist := user["uid"]; exist && a.Config.Core.UserStore != nil {
			userBase, err := a.Config.Core.UserStore.Get(uid, c)
			if err != nil {
				DefaultErrorResponse(c, errors.New("User not found!"))
				c.Abort()
				return
			}
			if !userBase.IsActive() {
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
