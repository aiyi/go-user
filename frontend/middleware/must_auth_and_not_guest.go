package middleware

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

// 检查客户端是否是认证状态, 并且不是 guest 认证.
// 如果是, 添加 sso_token_string<-->x-token, sso_token<-->*token.Token, sso_session<-->*session.Session, sso_user<-->*model.User  到 ctx *gin.Context;
// 如果否, 终止 Handlers Chain.
func MustAuthAndNotGuestHandler(ctx *gin.Context) {
	tkString := ctx.Request.Header.Get("x-token")
	if tkString == "" {
		ctx.JSON(200, errors.ErrTokenMissing)
		ctx.Abort()
		return
	}

	// 解析 token 并判断 AuthType 是否正确, 是否过期
	tk := &token.Token{}
	if err := tk.Decode([]byte(tkString)); err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrTokenDecodeFailed)
		ctx.Abort()
		return
	}
	if tk.AuthType == token.AuthTypeGuest {
		ctx.JSON(200, errors.ErrTokenShouldNotGuest)
		ctx.Abort()
		return
	}
	if time.Now().Unix() >= tk.ExpirationAccess {
		ctx.JSON(200, errors.ErrTokenAccessExpired)
		ctx.Abort()
		return
	}

	// 获取 Session 并判断与 token 是否匹配
	ss, err := session.Get(tk.SessionId)
	if err != nil {
		glog.Errorln(err)
		if err == errors.ErrNotFound {
			ctx.JSON(200, errors.ErrTokenInvalid)
			ctx.Abort()
			return
		}
		ctx.JSON(200, errors.ErrInternalServerError)
		ctx.Abort()
		return
	}
	if !security.SecureCompareString(tk.Signatrue, ss.TokenSignature) {
		ctx.JSON(200, errors.ErrTokenInvalid)
		ctx.Abort()
		return
	}

	// 获取用户信息并判断与 token, session 是否一致
	user, err := model.Get(ss.UserId)
	if err != nil {
		glog.Errorln(err)
		if err == model.ErrNotFound {
			ctx.JSON(200, errors.ErrTokenInvalid)
			ctx.Abort()
			return
		}
		ctx.JSON(200, errors.ErrInternalServerError)
		ctx.Abort()
		return
	}
	bindType, err := token.GetBindType(tk.AuthType)
	if err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrTokenInvalid)
		ctx.Abort()
		return
	}
	if user.BindTypes&bindType == 0 {
		ctx.JSON(200, errors.ErrTokenInvalid)
		ctx.Abort()
		return
	}
	if tk.AuthType == token.AuthTypeEmailPassword || tk.AuthType == token.AuthTypePhonePassword {
		if ss.PasswordTag != user.PasswordTag {
			ctx.JSON(200, errors.ErrTokenInvalid)
			ctx.Abort()
			return
		}
	}

	ctx.Set("sso_token_string", tkString)
	ctx.Set("sso_token", tk)
	ctx.Set("sso_session", ss)
	ctx.Set("sso_user", user)
}
