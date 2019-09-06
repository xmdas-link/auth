package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
)

func DefaultLogoutHandler(c *gin.Context, p AuthProvider, r AuthRender) {

	var (
		token = c.GetHeader("X-Token")
		gAuth  = GetAuthFromContext(c)
		err   error
	)

	if token == "" {
		err = errors.New("Token is required!")
	} else if gAuth == nil {
		err = errors.New("GinAuth missing!")
	} else if dbErr := gAuth.Config.Core.AuthToken.ClearToken(token); dbErr != nil {
		err = errors.New("清除token发生错误！")
	}

	if err != nil {
		c.Set("err", err.Error())
	}

	p.OnLogout(c)

	ret, err := r.Logout(c)
	doResponse(c, ret, err)
}
