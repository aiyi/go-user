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
	tk2EncodedBytes, err := tk2.Encode()
	if err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrTokenEncodeFailed)
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
