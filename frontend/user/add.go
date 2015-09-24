package user

import (
	"github.com/gin-gonic/gin"

	"github.com/aiyi/go-user/frontend"
)

func init() {
	frontend.UserRouter.POST("", AddHandler)
}

func AddHandler(ctx *gin.Context) {
	ctx.String(200, "ok")
}
