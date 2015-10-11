package frontend

import (
	"github.com/gin-gonic/gin"

	"github.com/aiyi/go-user/frontend/checkcode"
	"github.com/aiyi/go-user/frontend/middleware"
	"github.com/aiyi/go-user/frontend/user"
)

var Engine *gin.Engine

func init() {
	Engine = gin.Default()

	UserGroupRouter := Engine.Group("/user")
	{
		UserGroupRouter.GET("/auth", user.AuthHandler)
	}

	CheckCodeGroupRouter := Engine.Group("/checkcode")
	{
		CheckCodeGroupRouter.POST("/request_for_phone", middleware.CheckTokenHandler, checkcode.RequestForPhoneHandler)
		CheckCodeGroupRouter.POST("/request_for_email", middleware.CheckTokenHandler, checkcode.RequestForEmailHandler)
	}
}
