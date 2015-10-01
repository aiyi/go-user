package frontend

import (
	"crypto/hmac"
	"crypto/sha1"
	"github.com/chanxuehong/util/random"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"

	"github.com/aiyi/go-user/model"
)

func init() {
	UserRouter.POST("", AddHandler)
}

func AddHandler(ctx *gin.Context) {
	switch ctx.Request.Header.Get("auth_type") {
	case "email_password":
		addByEmailPassword(ctx)
	case "phone":
	case "qq":
	case "wechat":
	case "weibo":
	default:
	}
}

func addByEmailPassword(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email"    binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(200, NewError(1000, err.Error()))
		return
	}

	salt := random.NewRandom()
	h := hmac.New(sha1.New, salt[:])
	h.Write([]byte(req.Password))
	password := h.Sum(nil)

	userId, err := model.AddByEmail(req.Email, "", password, salt[:], time.Now().Unix())
	if err != nil {
		ctx.JSON(200, NewError(1000, err.Error()))
		return
	}

	ctx.JSON(200, NewError(200, strconv.FormatInt(userId, 10)))
}
