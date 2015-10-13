package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"github.com/aiyi/go-user/frontend/errors"
	"github.com/aiyi/go-user/frontend/session"
	"github.com/aiyi/go-user/frontend/token"
)

// 刷新 token
func RefreshHandler(ctx *gin.Context) {
	// MustAuthHandler(ctx)
	tk := ctx.MustGet("token").(*token.Token)
	if tk.AuthType == token.AuthTypeGuest { // 如果是 guest 用户直接返回
		resp := struct {
			*errors.Error
			Token string `json:"token"`
		}{
			Error: errors.ErrOK,
			Token: ctx.MustGet("token_string").(string),
		}
		ctx.JSON(200, &resp)
		return
	}

	timestamp := time.Now().Unix()
	if timestamp >= tk.ExpirationRefresh {
		ctx.JSON(200, errors.ErrTokenRefreshExpired)
		return
	}
	if timestamp+1200 < tk.ExpirationAccess { // 过早的刷新也是直接返回
		resp := struct {
			*errors.Error
			Token string `json:"token"`
		}{
			Error: errors.ErrOK,
			Token: ctx.MustGet("token_string").(string),
		}
		ctx.JSON(200, &resp)
		return
	}
	if tk.ExpirationAccess >= tk.ExpirationRefresh { // 暴力的结束, 防止客户端循环的刷新
		ctx.JSON(200, errors.ErrTokenRefreshExpired)
		return
	}

	tk2 := token.Token{
		SessionId:         tk.SessionId,
		TokenId:           token.NewTokenId(),
		AuthType:          tk.AuthType,
		ExpirationAccess:  token.ExpirationAccess(timestamp),
		ExpirationRefresh: tk.ExpirationRefresh,
	}
	if tk2.ExpirationAccess > tk2.ExpirationRefresh {
		tk2.ExpirationAccess = tk2.ExpirationRefresh
	}
	tk2EncodedBytes, err := tk2.Encode()
	if err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrTokenEncodeFailed)
		return
	}

	ss := ctx.MustGet("session").(*session.Session)
	ss.TokenSignature = tk2.Signatrue
	if err = session.Set(tk2.SessionId, ss); err != nil {
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
