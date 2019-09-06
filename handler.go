package auth

import "github.com/gin-gonic/gin"

type GinAuthHandler func(c *gin.Context, p AuthProvider, r AuthRender)

type GinAuthHandlers struct {
	GetLoginHandler      GinAuthHandler
	PostLoginHandler     GinAuthHandler
	LogoutHandler        GinAuthHandler
	LoginCallbackHandler GinAuthHandler
}

func NewGinHandler(h GinAuthHandler, a *GinAuth) func(c *gin.Context) {
	return func(c *gin.Context) {

		var (
			pName = c.Param("name")
			authP = a.GetProvider(pName)
			authR = a.GetRender(pName)
		)

		c.Set(CtxKeyGinAuth, a)
		c.Set(CtxKeyResultType, a.Config.Response.DefaultResultType)

		if authP == nil || authR == nil {
			if a.Config.Response.NotFoundHandler != nil {
				a.Config.Response.NotFoundHandler(c)
			} else {
				DefaultNotFoundHandler(c)
			}
			return
		}

		h(c, authP, authR)
	}
}
