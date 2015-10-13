package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"github.com/aiyi/go-user/frontend/errors"
	"github.com/aiyi/go-user/frontend/session"
	"github.com/aiyi/go-user/frontend/token"
	"github.com/aiyi/go-user/securitykey"
)

// 检查客户端是否是认证状态,
// 如果是, 添加 token<-->*token.Token, session<-->*session.Session  到 ctx *gin.Context,
// 如果不是, 终止 Handlers Chain.
func MustAuthHandler(ctx *gin.Context) {
	tkString := ctx.Request.Header.Get("x-token")
	if tkString == "" {
		ctx.JSON(200, errors.ErrTokenMissing)
		ctx.Abort()
		return
	}

	tk := &token.Token{}
	if err := tk.Decode([]byte(tkString), securitykey.Key); err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrTokenDecode)
		ctx.Abort()
		return
	}
	if tk.AuthType != token.AuthTypeGuest && time.Now().Unix() >= tk.ExpirationAccess {
		ctx.JSON(200, errors.ErrTokenAccessExpired)
		ctx.Abort()
		return
	}

	ss, err := session.Get(tk.SessionId)
	if err != nil {
		if err == errors.ErrNotFound {
			glog.Errorln(err)
			ctx.JSON(200, errors.ErrNotAuthOrExpired)
			ctx.Abort()
			return
		}
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrInternalServerError)
		ctx.Abort()
		return
	}
	if ss.TokenSignature != tk.Signatrue {
		ctx.JSON(200, errors.ErrNotAuthOrExpired)
		ctx.Abort()
		return
	}

	ctx.Set("token_string", tkString)
	ctx.Set("token", tk)
	ctx.Set("session", ss)
}
