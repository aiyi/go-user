package main

import (
	"fmt"

	"github.com/aiyi/go-user/model"
)

func add() {
	email := model.AddEmailParams{
		Email:    "email1@xxx.com",
		Password: []byte("password"),
		Salt:     []byte("salt"),
	}
	fmt.Println(model.AddEmail(&email))

	phone := model.AddPhoneParams{
		Phone:    "18888888888",
		Password: []byte("password"),
		Salt:     []byte("salt"),
	}
	fmt.Println(model.AddPhone(&phone))

	qq := model.AddQQParams{
		OpenId:   "openid",
		Nickname: "nickname",
	}
	fmt.Println(model.AddQQ(&qq))

	wechat := model.AddWechatParams{
		OpenId:   "openid",
		Nickname: "nickname",
	}
	fmt.Println(model.AddWechat(&wechat))

	weibo := model.AddWeiboParams{
		OpenId:   "openid",
		Nickname: "nickname",
	}
	fmt.Println(model.AddWeibo(&weibo))
}
