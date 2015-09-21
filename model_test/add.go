package main

import (
	"fmt"
	"log"

	"github.com/aiyi/go-user/model"
)

func add() {
	//email	 757120565977810729
	//phone	 757120566451767081
	//qq	 757120566976055081
	//wechat 757120567823304489
	//weibo	 757120568624416553
	emailUserId, err := model.AddEmail("email1@xxx.com", []byte("password"), []byte("salt"), 0)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("email\t", emailUserId)

	phoneUserId, err := model.AddPhone("18888888888", []byte("password"), []byte("salt"), 0)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("phone\t", phoneUserId)

	qqUserId, err := model.AddQQ("openid", "nickname", 0)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("qq\t", qqUserId)

	wechatUserId, err := model.AddWechat("openid", "nickname", 0)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("wechat\t", wechatUserId)

	weiboUserId, err := model.AddWeibo("openid", "nickname", 0)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("weibo\t", weiboUserId)
}
