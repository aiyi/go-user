package checkcode

import (
	"github.com/chanxuehong/util/check"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"github.com/aiyi/go-user/frontend/errors"
	"github.com/aiyi/go-user/frontend/session"
	"github.com/aiyi/go-user/frontend/token"
)

func RequestForPhoneHandler(ctx *gin.Context) {
	// CheckTokenHandler(ctx)
	queryValues := ctx.Request.URL.Query()
	phone := queryValues.Get("phone")
	if phone == "" {
		ctx.JSON(200, errors.ErrBadRequest)
		return
	}
	if !check.IsChinaMobileString(phone) {
		ctx.JSON(200, errors.ErrBadRequest)
		return
	}

	code := generateCode()
	if err := sendToPhone(phone, code); err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrInternalServerError)
		return
	}

	tk := ctx.MustGet("token").(*token.Token)
	ss := ctx.MustGet("session").(*session.Session)

	checkcode := session.CheckCode{
		Key:   phone,
		Code:  code,
		Times: 0,
	}
	ss.PhoneCheckCode = &checkcode
	if err := session.Set(tk.SessionId, ss); err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrInternalServerError)
		return
	}

	ctx.JSON(200, errors.ErrOK)
	return
}

func RequestForEmailHandler(ctx *gin.Context) {

}
