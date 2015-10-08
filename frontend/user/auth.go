package user

import (
	"time"

	"github.com/chanxuehong/util/security"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"github.com/aiyi/go-user/frontend/errors"
	"github.com/aiyi/go-user/frontend/sessiontoken"
	"github.com/aiyi/go-user/model"
	"github.com/aiyi/go-user/securitykey"
)

func AuthHandler(ctx *gin.Context) {
	switch authType := ctx.Request.Header.Get("x-auth-type"); authType {
	case sessiontoken.AuthTypeGuest:
		authGuestHandler(ctx)
	case sessiontoken.AuthTypeEmailPassword:
		authEmailPasswordHandler(ctx)
	case sessiontoken.AuthTypeEmailCheckcode:
	case sessiontoken.AuthTypePhonePassword:
		authPhonePasswordHandler(ctx)
	case sessiontoken.AuthTypePhoneCheckcode:
	case sessiontoken.AuthTypeOAuthQQ:
	case sessiontoken.AuthTypeOAuthWechat:
	case sessiontoken.AuthTypeOAuthWeibo:
	case "":
		ctx.JSON(200, errors.ErrAuthTypeMissing)
	default:
		ctx.JSON(200, errors.ErrAuthTypeUnknown)
	}
}

func authGuestHandler(ctx *gin.Context) {
	sid, err := sessiontoken.NewGuestSessionId()
	if err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrInternalServerError)
		return
	}

	timestamp := time.Now().Unix()

	tk := sessiontoken.SessionToken{
		SessionId:         sid,
		TokenId:           sessiontoken.NewTokenId(),
		AuthType:          sessiontoken.AuthTypeGuest,
		ExpirationAccess:  sessiontoken.ExpirationAccess(timestamp),
		ExpirationRefresh: sessiontoken.ExpirationRefresh(timestamp),
	}

	tkEncodedBytes, err := tk.Encode(securitykey.Key)
	if err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrSessionTokenEncode)
		return
	}

	ss := sessiontoken.Session{
		SessionTokenSignature: tk.Signatrue,
	}
	if err = sessiontoken.SessionAdd(sid, &ss); err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrInternalServerError)
		return
	}

	resp := struct {
		*errors.Error
		SessionToken string `json:"token"`
	}{
		Error:        errors.ErrOK,
		SessionToken: string(tkEncodedBytes),
	}
	ctx.JSON(200, &resp)
	return
}

func authSuccessHandler(ctx *gin.Context, authType string, user *model.User) {
	sid, err := sessiontoken.NewSessionId()
	if err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrInternalServerError)
		return
	}

	timestamp := time.Now().Unix()

	tk := sessiontoken.SessionToken{
		SessionId:         sid,
		TokenId:           sessiontoken.NewTokenId(),
		AuthType:          authType,
		ExpirationAccess:  sessiontoken.ExpirationAccess(timestamp),
		ExpirationRefresh: sessiontoken.ExpirationRefresh(timestamp),
	}

	tkEncodedBytes, err := tk.Encode(securitykey.Key)
	if err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrSessionTokenEncode)
		return
	}

	ss := sessiontoken.Session{
		SessionTokenSignature: tk.Signatrue,
		UserId:                user.Id,
		PasswordTag:           user.PasswordTag,
	}
	if err = sessiontoken.SessionAdd(sid, &ss); err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrInternalServerError)
		return
	}

	resp := struct {
		*errors.Error
		SessionToken string `json:"token"`
	}{
		Error:        errors.ErrOK,
		SessionToken: string(tkEncodedBytes),
	}
	ctx.JSON(200, &resp)
	return
}

func authEmailPasswordHandler(ctx *gin.Context) {
	querylValues := ctx.Request.URL.Query()
	email := querylValues.Get("email")
	if email == "" {
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

	// 匹配成功
	authSuccessHandler(ctx, sessiontoken.AuthTypeEmailPassword, user)
}

func authEmailCheckcode(ctx *gin.Context) {

}

func authPhonePasswordHandler(ctx *gin.Context) {
	querylValues := ctx.Request.URL.Query()
	phone := querylValues.Get("phone")
	if phone == "" {
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

	// 匹配成功
	authSuccessHandler(ctx, sessiontoken.AuthTypePhonePassword, user)
}
