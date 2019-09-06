package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
)

type Config struct {
	Path struct {
		// Auth工具挂载的地址。如：/auth
		Mount string
		// Auth工具挂载的Web服务的站点基础路径。如：http://localhost:8090/api，主要用于构建回调地址用。
		RootURL string
		// 域名
		Domain string
		// 登录成功之后的跳转地址，如果需要的话
		RedirectAfterLogin string
		// Token写入cookie时使用名字
		TokenKey string
	}

	Response struct {
		// 默认返回类型
		DefaultResultType ResultType

		// TODO：模板加载方式需要重新设计
		// TemplateFolder string

		NotFoundHandler func(c *gin.Context)
		// 捕获到预设之外的错误时，显示错误用的
		ErrorHandler func(c *gin.Context)
	}

	Core struct {

		// 路由
		Router *gin.Engine

		// 数据库连接
		DB *gorm.DB

		// token管理
		AuthToken AuthToken

		// Session
		Session Session

		// 是否打开日志功能
		LogMod bool
		Logger log.Logger
	}
}
