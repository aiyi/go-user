package web

import (
	"github.com/chanxuehong/wechat/open/oauth2"

	"github.com/aiyi/go-user/config"
)

var oauth2Config = oauth2.NewConfig(
	config.ConfigData.Weixin.Open.Web.AppId,
	config.ConfigData.Weixin.Open.Web.AppSecret,
	"unused",
	"snsapi_login",
)
