package frontend

import (
	"github.com/gin-gonic/gin"

	"github.com/aiyi/go-user/frontend/checkcode"
	"github.com/aiyi/go-user/frontend/user"
)

var (
	Engine = gin.Default()

	UserGroupRouter      = Engine.Group("/user")
	TokenGroupRouter     = Engine.Group("/token")
	CheckCodeGroupRouter = Engine.Group("/checkcode")
)

func init() {
	UserGroupRouter.GET("/auth", user.AuthHandler)
	CheckCodeGroupRouter.POST("/request_code", checkcode.GetCheckCodeHandler)
}
