package frontend

import (
	"github.com/gin-gonic/gin"

	"github.com/aiyi/go-user/frontend/checkcode"
	"github.com/aiyi/go-user/frontend/middleware"
	"github.com/aiyi/go-user/frontend/oauth/wechat/mp"
	"github.com/aiyi/go-user/frontend/user"
)

var Engine *gin.Engine

func init() {
	Engine = gin.Default()

	Engine.GET("/oauth/wechat/mp/login_url", middleware.MustAuthHandler, mp.LoginURLHandler)

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
