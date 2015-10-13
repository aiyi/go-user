package handler

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

// 刷新 token
func RefreshHandler(ctx *gin.Context) {
	tkString := ctx.Request.Header.Get("x-token")
	if tkString == "" {
		ctx.JSON(200, errors.ErrTokenMissing)
		return
	}

	// 解析 token
	tk := &token.Token{}
	if err := tk.Decode([]byte(tkString)); err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrTokenDecodeFailed)
		return
	}

	// 如果是 guest 直接返回
	if tk.AuthType == token.AuthTypeGuest {
		resp := struct {
			*errors.Error
			Token string `json:"token"`
		}{
			Error: errors.ErrOK,
			Token: tkString,
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
			Token: tkString,
		}
		ctx.JSON(200, &resp)
		return
	}
	if tk.ExpirationAccess >= tk.ExpirationRefresh { // 提前 1200秒 暴力的结束, 防止客户端循环的刷新
		ctx.JSON(200, errors.ErrTokenRefreshExpired)
		return
	}

	// 现在我们真的需要刷新token啦!!!

	// 获取 Session 并判断与 token 是否匹配
	ss, err := session.Get(tk.SessionId)
	if err != nil {
		glog.Errorln(err)
		if err == errors.ErrNotFound {
			ctx.JSON(200, errors.ErrTokenInvalid)
			return
		}
		ctx.JSON(200, errors.ErrInternalServerError)
		return
	}
	if !security.SecureCompareString(tk.Signatrue, ss.TokenSignature) {
		ctx.JSON(200, errors.ErrTokenInvalid)
		return
	}

	// 获取用户信息并判断与 token, session 是否一致
	user, err := model.Get(ss.UserId)
	if err != nil {
		glog.Errorln(err)
		if err == model.ErrNotFound {
			ctx.JSON(200, errors.ErrTokenInvalid)
			return
		}
		ctx.JSON(200, errors.ErrInternalServerError)
		return
	}
	bindType, err := token.GetBindType(tk.AuthType)
	if err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrTokenInvalid)
		return
	}
	if user.BindTypes&bindType == 0 {
		ctx.JSON(200, errors.ErrTokenInvalid)
		return
	}
	if tk.AuthType == token.AuthTypeEmailPassword || tk.AuthType == token.AuthTypePhonePassword {
		if ss.PasswordTag != user.PasswordTag {
			ctx.JSON(200, errors.ErrTokenInvalid)
			return
		}
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
