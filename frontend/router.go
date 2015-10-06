package frontend

import (
	"github.com/gin-gonic/gin"
)

var (
	Engine = gin.Default()

	UserGroupRouter  = Engine.Group("/user")
	TokenGroupRouter = Engine.Group("/token")
)
