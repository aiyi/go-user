package user

import (
	"net/url"
	"time"

	"github.com/chanxuehong/util/check"
	"github.com/chanxuehong/util/security"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"github.com/aiyi/go-user/frontend/errors"
	"github.com/aiyi/go-user/frontend/session"
	"github.com/aiyi/go-user/frontend/token"
	"github.com/aiyi/go-user/model"
)

// 认证, 获取token.
//  uri?auth_type=guest
//  uri?auth_type=email_password&email=XXX&password=XXX
//  uri?auth_type=phone_password&phone=XXX&password=XXX
//  uri?auth_type=email_checkcode&email=XXX&checkcode=XXX
//  uri?auth_type=phone_checkcode&phone=XXX&checkcode=XXX
func AuthHandler(ctx *gin.Context) {
	queryValues := ctx.Request.URL.Query()
	authType := queryValues.Get("auth_type")
	switch authType {
	case AuthTypeGuest:
		authGuest(ctx, queryValues)
	case AuthTypeEmailPassword:
		authEmailPassword(ctx, queryValues)
	case AuthTypeEmailCheckCode:
		ctx.JSON(200, errors.ErrNotSupported)
	case AuthTypePhonePassword:
		authPhonePassword(ctx, queryValues)
	case AuthTypePhoneCheckCode:
		ctx.JSON(200, errors.ErrNotSupported)
	default:
		ctx.JSON(200, errors.ErrBadRequest)
	}
}

func authGuest(ctx *gin.Context, queryValues url.Values) {
	sid, err := session.NewGuestSessionId()
	if err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrInternalServerError)
		return
	}

	tk := token.Token{
		SessionId:         sid,
		TokenId:           token.NewTokenId(),
		AuthType:          AuthTypeGuest,
		ExpirationAccess:  0,
		ExpirationRefresh: 0,
	}
	tkEncodedBytes, err := tk.Encode()
	if err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrTokenEncodeFailed)
		return
	}

	ss := session.Session{
		TokenSignature: tk.Signatrue,
	}
	if err = session.Add(sid, &ss); err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrInternalServerError)
		return
	}

	resp := struct {
		*errors.Error
		Token string `json:"token"`
	}{
		Error: errors.ErrOK,
		Token: string(tkEncodedBytes),
	}
	ctx.JSON(200, &resp)
	return
}

// 认证成功后创建 token 和 session
func authSuccess(ctx *gin.Context, authType string, user *model.User) {
	sid, err := session.NewSessionId()
	if err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrInternalServerError)
		return
	}

	timestamp := time.Now().Unix()
	tk := token.Token{
		SessionId:         sid,
		TokenId:           token.NewTokenId(),
		AuthType:          authType,
		ExpirationAccess:  token.ExpirationAccess(timestamp),
		ExpirationRefresh: token.ExpirationRefresh(timestamp),
	}
	tkEncodedBytes, err := tk.Encode()
	if err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrTokenEncodeFailed)
		return
	}

	ss := session.Session{
		TokenSignature: tk.Signatrue,
		UserId:         user.Id,
		PasswordTag:    user.PasswordTag,
	}
	if err = session.Add(sid, &ss); err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrInternalServerError)
		return
	}

	resp := struct {
		*errors.Error
		Token string `json:"token"`
	}{
		Error: errors.ErrOK,
		Token: string(tkEncodedBytes),
	}
	ctx.JSON(200, &resp)
	return
}

func authEmailPassword(ctx *gin.Context, querylValues url.Values) {
	email := querylValues.Get("email")
	if email == "" {
		ctx.JSON(200, errors.ErrBadRequest)
		return
	}
	if !check.IsMailString(email) {
		ctx.JSON(200, errors.ErrBadRequest)
		return
	}
	password := querylValues.Get("password")
	if password == "" {
		ctx.JSON(200, errors.ErrBadRequest)
		return
	}

	user, err := model.GetByEmail(email)
	switch err {
	case nil:
		cipherPassword := model.EncryptPassword([]byte(password), user.Salt)
		if !security.SecureCompare(cipherPassword, user.Password) {
			ctx.JSON(200, errors.ErrAuthFailed)
			return
		}
		authSuccess(ctx, AuthTypeEmailPassword, user)
		return
	case model.ErrNotFound:
		cipherPassword := model.EncryptPassword([]byte(password), model.PasswordSalt)
		if !security.SecureCompare(cipherPassword, cipherPassword) {
			ctx.JSON(200, errors.ErrAuthFailed)
			return
		}
		ctx.JSON(200, errors.ErrAuthFailed)
		return
	default:
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrInternalServerError)
		return
	}
}

func authPhonePassword(ctx *gin.Context, querylValues url.Values) {
	phone := querylValues.Get("phone")
	if phone == "" {
		ctx.JSON(200, errors.ErrBadRequest)
		return
	}
	if !check.IsChinaMobileString(phone) {
		ctx.JSON(200, errors.ErrBadRequest)
		return
	}
	password := querylValues.Get("password")
	if password == "" {
		ctx.JSON(200, errors.ErrBadRequest)
		return
	}

	user, err := model.GetByPhone(phone)
	switch err {
	case nil:
		cipherPassword := model.EncryptPassword([]byte(password), user.Salt)
		if !security.SecureCompare(cipherPassword, user.Password) {
			ctx.JSON(200, errors.ErrAuthFailed)
			return
		}
		authSuccess(ctx, AuthTypePhonePassword, user)
		return
	case model.ErrNotFound:
		cipherPassword := model.EncryptPassword([]byte(password), model.PasswordSalt)
		if !security.SecureCompare(cipherPassword, cipherPassword) {
			ctx.JSON(200, errors.ErrAuthFailed)
			return
		}
		ctx.JSON(200, errors.ErrAuthFailed)
		return
	default:
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrInternalServerError)
		return
	}
}
