package sessiontoken

import (
	"crypto/hmac"
	"crypto/sha1"
	"strconv"
	"time"

	"github.com/chanxuehong/util/random"
	"github.com/gin-gonic/gin"

	"github.com/aiyi/go-user/frontend/errors"
	"github.com/aiyi/go-user/model"
)

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
		ctx.JSON(200, errors.NewError(1000, err.Error()))
		return
	}

	salt := random.NewRandom()
	h := hmac.New(sha1.New, salt[:])
	h.Write([]byte(req.Password))
	password := h.Sum(nil)

	userId, err := model.AddByEmail(req.Email, "", password, salt[:], time.Now().Unix())
	if err != nil {
		ctx.JSON(200, errors.NewError(1000, err.Error()))
		return
	}

	ctx.JSON(200, errors.NewError(200, strconv.FormatInt(userId, 10)))
}
