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

		c.Set(CtxKeyUserRole, user["role"])
		c.Set(CtxKeyAuthUser, user)

		c.Next()
	}

}
