package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
)

func DefaultLogoutHandler(c *gin.Context, gAuth *GinAuth) {

	var (
		authToken = gAuth.Config.Core.AuthToken
		token     string
		err       error
	)

	if gAuth == nil {
		err = errors.New("GinAuth missing!")
	} else if authToken == nil {
		err = errors.New("AuthToken missing!")
	} else {
		token = c.GetHeader(gAuth.Config.Path.HeaderKey)
		if token == "" {
			err = errors.New("Token is required!")
		} else if user := authToken.FindToken(token); user == nil {
			// token不存在，也就没有登出的问题
			DefaultJsonSuccess(c, nil)
			return
		} else if dbErr := gAuth.Config.Core.AuthToken.ClearToken(token); dbErr != nil {
			err = errors.New("清除token发生错误！")
		} else {
			providerName := user["provider"]
			p := gAuth.GetProvider(providerName)
			r := gAuth.GetRender(providerName)
			if p == nil || r == nil {
				// err = errors.New("Provider or Render missing!")
				// 不知道渲染方式的话，返回JSON
				DefaultJsonSuccess(c, nil)
			} else {
				p.OnLogout(c)
				ret, err := r.Logout(c)
				doResponse(c, ret, err)
				return
			}

		}
	}

	if err != nil {
		DefaultJsonError(c, err.Error())
	}
}
