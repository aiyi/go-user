package open

import (
	"time"

	"github.com/chanxuehong/util/security"
	"github.com/chanxuehong/wechat/oauth2"
	openOAuth2 "github.com/chanxuehong/wechat/open/oauth2"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"github.com/aiyi/go-user/frontend/errors"
	"github.com/aiyi/go-user/frontend/session"
	"github.com/aiyi/go-user/frontend/token"
	"github.com/aiyi/go-user/model"
)

// 微信 oauth2 登录
//  需要提供 code, state 参数.
func LoginHandler(ctx *gin.Context) {
	// MustAuthHandler(ctx)
	queryValues := ctx.Request.URL.Query()
	code := queryValues.Get("code")
	if code == "" {
		ctx.JSON(200, errors.ErrBadRequest)
		return
	}
	state := queryValues.Get("state")
	if state == "" {
		ctx.JSON(200, errors.ErrBadRequest)
		return
	}

	ss := ctx.MustGet("sso_session").(*session.Session)

	// 比较 state 是否一致
	if !security.SecureCompareString(state, ss.OAuth2State) {
		ctx.JSON(200, errors.ErrOAuthStateMismatch)
		return
	}

	oauth2Client := oauth2.Client{
		Config: oauth2Config,
	}
	oauth2Token, err := oauth2Client.Exchange(code)
	if err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrInternalServerError)
		return
	}
	if oauth2Token.UnionId == "" {
		glog.Errorln("unionid is empty")
		ctx.JSON(200, errors.ErrInternalServerError)
		return
	}

	var oauth2UserInfo openOAuth2.UserInfo
	if err = oauth2Client.GetUserInfo(&oauth2UserInfo, ""); err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrInternalServerError)
		return
	}

	timestamp := time.Now().Unix()
	user, err := model.GetByWechat(oauth2Token.UnionId)
	switch err {
	default:
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrInternalServerError)
		return
	case model.ErrNotFound:
		user, err = model.AddByWechat(oauth2Token.UnionId, oauth2UserInfo.Nickname, timestamp)
		if err != nil {
			glog.Errorln(err)
			ctx.JSON(200, errors.ErrInternalServerError)
			return
		}
		fallthrough
	case nil:
		sid, err := session.NewSessionId()
		if err != nil {
			glog.Errorln(err)
			ctx.JSON(200, errors.ErrInternalServerError)
			return
		}

		tk2 := token.Token{
			SessionId:         sid,
			TokenId:           token.NewTokenId(),
			AuthType:          token.AuthTypeOAuthWechat,
			ExpirationAccess:  token.ExpirationAccess(timestamp),
			ExpirationRefresh: token.ExpirationRefresh(timestamp),
		}
		tk2EncodedBytes, err := tk2.Encode()
		if err != nil {
			glog.Errorln(err)
			ctx.JSON(200, errors.ErrTokenEncodeFailed)
			return
		}

		ss2 := session.Session{
			TokenSignature: tk2.Signatrue,
			UserId:         user.Id,
			PasswordTag:    user.PasswordTag,
		}
		if err = session.Add(sid, &ss2); err != nil {
			glog.Errorln(err)
			ctx.JSON(200, errors.ErrInternalServerError)
			return
		}

		resp := struct {
			*errors.Error
			Token string `json:"token"`
		}{
			Error: errors.ErrOK,
			Token: string(tk2EncodedBytes),
		}
		ctx.JSON(200, &resp)
		return
	}
}
