package main

import (
	"fmt"
	"time"

	"github.com/aiyi/go-user/model"
	"github.com/aiyi/go-user/userid"
)

func main() {
	//	userId, err := userid.GetId()
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//	para := model.AddEmailParams{
	//		UserId:     userId,
	//		Email:      "test@test.test",
	//		Password:   []byte("password"),
	//		Salt:       []byte("salt"),
	//		CreateTime: 13799999,
	//	}
	//	fmt.Println(model.AddEmail(&para))

	//	userId, err := userid.GetId()
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//	para := model.AddPhoneParams{
	//		UserId:     userId,
	//		Phone:      "18988889999",
	//		Password:   []byte("password"),
	//		Salt:       []byte("salt"),
	//		CreateTime: time.Now().Unix(),
	//	}
	//	fmt.Println(model.AddPhone(&para))

	//	userId, err := userid.GetId()
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//	para := model.AddQQParams{
	//		UserId:     userId,
	//		OpenId:     "openid",
	//		CreateTime: time.Now().Unix(),
	//	}
	//	fmt.Println(model.AddQQ(&para))

	//	userId, err := userid.GetId()
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//	para := model.AddWechatParams{
	//		UserId:     userId,
	//		OpenId:     "openid",
	//		CreateTime: time.Now().Unix(),
	//	}
	//	fmt.Println(model.AddWechat(&para))

	userId, err := userid.GetId()
	if err != nil {
		fmt.Println(err)
		return
	}
	para := model.AddWeiboParams{
		UserId:     userId,
		OpenId:     "openid",
		CreateTime: time.Now().Unix(),
	}
	fmt.Println(model.AddWeibo(&para))
}
