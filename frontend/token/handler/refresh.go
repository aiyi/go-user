package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"github.com/aiyi/go-user/frontend/errors"
	"github.com/aiyi/go-user/frontend/session"
	"github.com/aiyi/go-user/frontend/token"
	"github.com/aiyi/go-user/securitykey"
)

// 刷新 token
func RefreshHandler(ctx *gin.Context) {
	// NOTE: 在此之前的中间件获取了 token 和 session
	tk := ctx.MustGet("token").(*token.Token)
	if tk.AuthType == token.AuthTypeGuest {
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

	tk2 := token.Token{
		SessionId:         tk.SessionId,
		TokenId:           token.NewTokenId(),
		AuthType:          tk.AuthType,
		ExpirationAccess:  token.ExpirationAccess(time.Now().Unix()),
		ExpirationRefresh: tk.ExpirationRefresh,
	}
	tk2EncodedBytes, err := tk2.Encode(securitykey.Key)
	if err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrTokenEncode)
		return
	}

	ss := ctx.MustGet("session").(*session.Session)
	ss.TokenSignature = tk2.Signatrue
	if err = session.Set(tk.SessionId, ss); err != nil {
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
