package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
)

func (a *GinAuth) AuthTokenMiddleware() func(c *gin.Context) {

	return func(c *gin.Context) {

		token := c.GetHeader("X-Token")
		// token = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Njc0ODk5MjYsIlVzZXIiOnsiSUQiOiIxIiwiSVAiOiIxMjcuMC4wLjEiLCJQcm92aWRlciI6InBhc3N3b3JkIiwiUm9sZSI6InRlc3QifX0.ai-iW63QmbxIPWwjYirk-y3pmSZNdVRpIyA4hWK9No3dOkT5Jd2BHlGNnyqScILN2t3a6OQbJfzG1NZRwx6W0KWF-BtqJTqL-TyInpj6QvJpN5o3t2k1_zu-QCRa-a3WZ-oC8xkgl_mgfWQp_YeQgnaH2mZm_O_PEasnJrW8FlQ"
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
