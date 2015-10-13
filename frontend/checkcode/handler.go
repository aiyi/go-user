package checkcode

import (
	"github.com/chanxuehong/util/check"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"github.com/aiyi/go-user/frontend/errors"
	"github.com/aiyi/go-user/frontend/session"
	"github.com/aiyi/go-user/frontend/token"
)

// 申请发送一个校验码到手机.
//  uri?phone=XXX
func RequestForPhoneHandler(ctx *gin.Context) {
	// MustAuthHandler(ctx)
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

	tk := ctx.MustGet("sso_token").(*token.Token)
	ss := ctx.MustGet("sso_session").(*session.Session)

	code := generateCode()
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

	if err := sendCodeToPhone(phone, code); err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrInternalServerError)
		return
	}

	ctx.JSON(200, errors.ErrOK)
	return
}

// 申请发送一个校验码到邮箱.
//  uri?email=XXX
func RequestForEmailHandler(ctx *gin.Context) {
	// MustAuthHandler(ctx)
	queryValues := ctx.Request.URL.Query()
	email := queryValues.Get("email")
	if email == "" {
		ctx.JSON(200, errors.ErrBadRequest)
		return
	}
	if !check.IsMailString(email) {
		ctx.JSON(200, errors.ErrBadRequest)
		return
	}

	tk := ctx.MustGet("sso_token").(*token.Token)
	ss := ctx.MustGet("sso_session").(*session.Session)

	code := generateCode()
	checkcode := session.CheckCode{
		Key:   email,
		Code:  code,
		Times: 0,
	}
	ss.EmailCheckCode = &checkcode
	if err := session.Set(tk.SessionId, ss); err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrInternalServerError)
		return
	}

	if err := sendCodeToEmail(email, code); err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrInternalServerError)
		return
	}

	ctx.JSON(200, errors.ErrOK)
	return
}
