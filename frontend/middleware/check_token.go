package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"github.com/aiyi/go-user/frontend/errors"
	"github.com/aiyi/go-user/frontend/session"
	"github.com/aiyi/go-user/frontend/token"
	"github.com/aiyi/go-user/securitykey"
)

// 检查客户端传递过来的 token 是否有效, 如果有效设置 *Token, *Session 到 gin.Context.
func CheckTokenHandler(ctx *gin.Context) {
	tkBytes := ctx.Request.Header.Get("x-token")
	if tkBytes == "" {
		ctx.JSON(200, errors.ErrTokenMissing)
		ctx.Abort()
		return
	}

	var tk token.Token
	if err := tk.Decode([]byte(tkBytes), securitykey.Key); err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrTokenDecode)
		ctx.Abort()
		return
	}

	ss, err := session.Get(tk.SessionId)
	if err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrInternalServerError)
		ctx.Abort()
		return
	}
	if ss.TokenSignature != tk.Signatrue {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrInternalServerError)
		ctx.Abort()
		return
	}

	ctx.Set("token", &tk)
	ctx.Set("session", ss)
}
