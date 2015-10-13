package middleware

import (
	"time"

	"github.com/chanxuehong/util/security"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"github.com/aiyi/go-user/frontend/errors"
	"github.com/aiyi/go-user/frontend/session"
	"github.com/aiyi/go-user/frontend/token"
)

// 检查客户端是否是认证状态, 并且不是 guest 用户.
// 如果是, 添加 token_string<-->x-token, token<-->*token.Token, session<-->*session.Session  到 ctx *gin.Context;
// 如果不是, 终止 Handlers Chain.
func MustNotGuestAuthHandler(ctx *gin.Context) {
	tkString := ctx.Request.Header.Get("x-token")
	if tkString == "" {
		ctx.JSON(200, errors.ErrTokenMissing)
		ctx.Abort()
		return
	}

	tk := &token.Token{}
	if err := tk.Decode([]byte(tkString)); err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrTokenDecode)
		ctx.Abort()
		return
	}
	if tk.AuthType == token.AuthTypeGuest {
		ctx.JSON(200, errors.ErrNotAuthOrExpired)
		ctx.Abort()
		return
	}
	if time.Now().Unix() >= tk.ExpirationAccess {
		ctx.JSON(200, errors.ErrTokenAccessExpired)
		ctx.Abort()
		return
	}

	ss, err := session.Get(tk.SessionId)
	if err != nil {
		glog.Errorln(err)
		if err == errors.ErrNotFound {
			ctx.JSON(200, errors.ErrNotAuthOrExpired)
			ctx.Abort()
			return
		}
		ctx.JSON(200, errors.ErrInternalServerError)
		ctx.Abort()
		return
	}
	if !security.SecureCompareString(tk.Signatrue, ss.TokenSignature) {
		ctx.JSON(200, errors.ErrNotAuthOrExpired)
		ctx.Abort()
		return
	}

	ctx.Set("token_string", tkString)
	ctx.Set("token", tk)
	ctx.Set("session", ss)
}
