package open

import (
	"github.com/chanxuehong/wechat/open/oauth2"

	"github.com/aiyi/go-user/config"
)

var oauth2Config = oauth2.NewConfig(
	config.ConfigData.Weixin.Open.AppId,
	config.ConfigData.Weixin.Open.AppSecret,
	"unused",
	"snsapi_login",
)
