package user

import (
	"time"

	"github.com/chanxuehong/util/security"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"github.com/aiyi/go-user/frontend/errors"
	"github.com/aiyi/go-user/frontend/session"
	"github.com/aiyi/go-user/frontend/token"
	"github.com/aiyi/go-user/model"
)

func AuthHandler(ctx *gin.Context) {
	switch authType := ctx.Request.Header.Get("x-auth-type"); authType {
	case AuthTypeGuest:
		authGuestHandler(ctx)
	case AuthTypeEmailPassword:
		authEmailPasswordHandler(ctx)
	case AuthTypeEmailCheckCode:
		authEmailCheckCodeHandler(ctx)
	case AuthTypePhonePassword:
		authPhonePasswordHandler(ctx)
	case AuthTypePhoneCheckCode:
	case AuthTypeOAuthQQ:
	case AuthTypeOAuthWechat:
	case AuthTypeOAuthWeibo:
	case "":
		ctx.JSON(200, errors.ErrAuthTypeMissing)
	default:
		ctx.JSON(200, errors.ErrAuthTypeUnknown)
	}
}

// 创建 guest token 和 session
func authGuestHandler(ctx *gin.Context) {
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
		ctx.JSON(200, errors.ErrTokenEncode)
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
func authSuccessHandler(ctx *gin.Context, authType string, user *model.User) {
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
		ctx.JSON(200, errors.ErrTokenEncode)
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
		authSuccessHandler(ctx, AuthTypeEmailPassword, user)
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
		ctx.JSON(200, errors.ErrTokenMissing)
		return
	}

	var tk token.Token
	if err := tk.Decode([]byte(tkBytes)); err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrTokenDecode)
		return
	}

	ss, err := session.Get(tk.SessionId)
	if err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrInternalServerError)
		return
	}
	if ss.TokenSignature != tk.Signatrue {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrInternalServerError)
		return
	}
	//	wantCheckCode := ss.EmailCheckCode
	//	if wantCheckCode == nil || time.Now().Unix() >= wantCheckCode.Expiration || wantCheckCode.Code != checkcode {
	//		glog.Errorln(err)
	//		ctx.JSON(200, errors.ErrInternalServerError)
	//		return
	//	}

	user, err := model.GetByEmail(email)
	switch err {
	case nil:
		authSuccessHandler(ctx, AuthTypeEmailCheckCode, user)
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
		authSuccessHandler(ctx, AuthTypePhonePassword, user)
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
