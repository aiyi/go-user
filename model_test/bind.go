package main

import (
	"fmt"

	"github.com/aiyi/go-user/model"
)

func bind() {
	email := model.BindEmailParams{
		UserId: 757050443187618548,
		Email:  "bind@xxx.com",
	}
	fmt.Println(model.BindEmail(&email))

	phone := model.BindPhoneParams{
		UserId: 757050443971953396,
		Phone:  "19999999999",
	}
	fmt.Println(model.BindPhone(&phone))

	qq := model.BindQQParams{
		UserId:   757050444483658484,
		OpenId:   "bind_openid",
		Nickname: "bind_nickname",
	}
	fmt.Println(model.BindQQ(&qq))

	wechat := model.BindWechatParams{
		UserId:   757050444982780660,
		OpenId:   "bind_openid",
		Nickname: "bind_nickname",
	}
	fmt.Println(model.BindWechat(&wechat))

	weibo := model.BindWeiboParams{
		UserId:   757048980382156277,
		OpenId:   "bind_openid",
		Nickname: "bind_nickname",
	}
	fmt.Println(model.BindWeibo(&weibo))
}
