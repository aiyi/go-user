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
	case sessiontoken.AuthTypeEmailCheckCode:
		authEmailCheckCodeHandler(ctx)
	case sessiontoken.AuthTypePhonePassword:
		authPhonePasswordHandler(ctx)
	case sessiontoken.AuthTypePhoneCheckCode:
	case sessiontoken.AuthTypeOAuthQQ:
	case sessiontoken.AuthTypeOAuthWechat:
	case sessiontoken.AuthTypeOAuthWeibo:
	case "":
		ctx.JSON(200, errors.ErrAuthTypeMissing)
	default:
		ctx.JSON(200, errors.ErrAuthTypeUnknown)
	}
}

// 创建 guest token 和 session
func authGuestHandler(ctx *gin.Context) {
	sid, err := sessiontoken.NewGuestSessionId()
	if err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrInternalServerError)
		return
	}

	tk := sessiontoken.SessionToken{
		SessionId:         sid,
		TokenId:           "",
		AuthType:          sessiontoken.AuthTypeGuest,
		ExpirationAccess:  0,
		ExpirationRefresh: 0,
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

// 认证成功后创建 token 和 session
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
		authSuccessHandler(ctx, sessiontoken.AuthTypeEmailPassword, user)
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

func authEmailCheckCodeHandler(ctx *gin.Context) {
	querylValues := ctx.Request.URL.Query()
	email := querylValues.Get("email")
	if email == "" {
		ctx.JSON(200, errors.ErrBadRequest)
		return
	}
	checkcode := querylValues.Get("checkcode")
	if checkcode == "" {
		ctx.JSON(200, errors.ErrBadRequest)
		return
	}

	tkBytes := ctx.Request.Header.Get("x-token")
	if tkBytes == "" {
		ctx.JSON(200, errors.ErrSessionTokenMissing)
		return
	}

	var tk sessiontoken.SessionToken
	if err := tk.Decode([]byte(tkBytes), securitykey.Key); err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrSessionTokenDecode)
		return
	}

	ss, err := sessiontoken.SessionGet(tk.SessionId)
	if err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrInternalServerError)
		return
	}
	if ss.SessionTokenSignature != tk.Signatrue {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrInternalServerError)
		return
	}
	wantCheckCode := ss.EmailCheckCode
	if wantCheckCode == nil || time.Now().Unix() >= wantCheckCode.Expiration || wantCheckCode.Code != checkcode {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrInternalServerError)
		return
	}

	user, err := model.GetByEmail(email)
	switch err {
	case nil:
		authSuccessHandler(ctx, sessiontoken.AuthTypeEmailCheckCode, user)
		return
	case model.ErrNotFound:
		ctx.JSON(200, errors.ErrAuthFailed)
		return
	default:
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrInternalServerError)
		return
	}
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
		authSuccessHandler(ctx, sessiontoken.AuthTypePhonePassword, user)
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
