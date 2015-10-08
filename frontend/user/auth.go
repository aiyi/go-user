package user

import (
	"time"

	"github.com/chanxuehong/util/random"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"github.com/aiyi/go-user/frontend/errors"
	"github.com/aiyi/go-user/frontend/sessiontoken"
	"github.com/aiyi/go-user/securitykey"
)

func AuthHandler(ctx *gin.Context) {
	switch AuthType := ctx.Request.Header.Get("auth_type"); AuthType {
	case sessiontoken.AuthTypeGuest:
		authGuest(ctx)
	case sessiontoken.AuthTypeEmailPassword:
		authEmailPassword(ctx)
	case sessiontoken.AuthTypeEmailCheckcode:
	case sessiontoken.AuthTypePhonePassword:
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

func authGuest(ctx *gin.Context) {
	sid, err := sessiontoken.NewGuestSessionId()
	if err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrInternalServerError)
		return
	}

	timestamp := time.Now().Unix()

	tk := sessiontoken.SessionToken{
		SessionId:         sid,
		TokenId:           string(random.NewRandomEx()),
		AuthType:          sessiontoken.AuthTypeGuest,
		ExpirationAccess:  sessiontoken.ExpirationAccess(timestamp),
		ExpirationRefresh: sessiontoken.ExpirationRefresh(timestamp),
	}

	tkBytes, err := tk.Encode(securitykey.Key)
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
		SessionToken string `json:"sesstiontoken"`
	}{
		Error:        errors.ErrOK,
		SessionToken: string(tkBytes),
	}
	ctx.JSON(200, &resp)
	return
}

func authEmailPassword(ctx *gin.Context) {

}

func authEmailCheckcode(ctx *gin.Context) {

}
