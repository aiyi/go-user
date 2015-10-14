package app

import (
	"github.com/chanxuehong/util/random"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"github.com/aiyi/go-user/config"
	"github.com/aiyi/go-user/frontend/errors"
	"github.com/aiyi/go-user/frontend/session"
	"github.com/aiyi/go-user/frontend/token"
)

// 获取请求用户授权的参数(appid, state, scope)
func AuthParaHandler(ctx *gin.Context) {
	// MustAuthHandler(ctx)
	tk := ctx.MustGet("sso_token").(*token.Token)
	ss := ctx.MustGet("sso_session").(*session.Session)

	ss.OAuth2State = string(random.NewRandomEx())
	if err := session.Set(tk.SessionId, ss); err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrInternalServerError)
		return
	}

	resp := struct {
		*errors.Error
		AppId string `json:"appid"`
		State string `json:"state"`
		Scope string `json:"scope"`
	}{
		Error: errors.ErrOK,
		AppId: config.ConfigData.Weixin.Open.App.AppId,
		State: ss.OAuth2State,
		Scope: "snsapi_userinfo",
	}
	ctx.JSON(200, &resp)
	return
}
