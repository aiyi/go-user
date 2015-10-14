package frontend

import (
	"github.com/gin-gonic/gin"

	"github.com/aiyi/go-user/frontend/checkcode"
	"github.com/aiyi/go-user/frontend/middleware"
	"github.com/aiyi/go-user/frontend/oauth/wechat/mp"
	"github.com/aiyi/go-user/frontend/oauth/wechat/open/app"
	"github.com/aiyi/go-user/frontend/oauth/wechat/open/web"
	tokenhandler "github.com/aiyi/go-user/frontend/token/handler"
	"github.com/aiyi/go-user/frontend/user"
)

var Engine *gin.Engine

func init() {
	Engine = gin.Default()

	Engine.GET("/oauth/wechat/mp/auth_url", middleware.MustAuthHandler, mp.AuthURLHandler)
	Engine.GET("/oauth/wechat/mp/auth", middleware.MustAuthHandler, mp.AuthHandler)
	Engine.GET("/oauth/wechat/open/web/auth_url", middleware.MustAuthHandler, web.AuthURLHandler)
	Engine.GET("/oauth/wechat/open/web/auth", middleware.MustAuthHandler, web.AuthHandler)
	Engine.GET("/oauth/wechat/open/app/auth_para", middleware.MustAuthHandler, app.AuthParaHandler)
	Engine.GET("/oauth/wechat/open/app/auth", middleware.MustAuthHandler, app.AuthHandler)

	Engine.GET("/token/refresh", tokenhandler.RefreshHandler)

	UserGroupRouter := Engine.Group("/user")
	{
		UserGroupRouter.GET("/auth", user.AuthHandler)
	}

	CheckCodeGroupRouter := Engine.Group("/checkcode")
	{
		CheckCodeGroupRouter.POST("/request_for_phone", middleware.MustAuthHandler, checkcode.RequestForPhoneHandler)
		CheckCodeGroupRouter.POST("/request_for_email", middleware.MustAuthHandler, checkcode.RequestForEmailHandler)
	}
}
