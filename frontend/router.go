package frontend

import (
	"github.com/gin-gonic/gin"
)

var (
	Engine      = gin.Default()
	UserRouter  = Engine.Group("/user")
	TokenRouter = UserRouter.Group("/token")
)
