package frontend

import (
	"github.com/gin-gonic/gin"
)

var (
	Router     = gin.Default()
	UserRouter = Router.Group("/user")
)
