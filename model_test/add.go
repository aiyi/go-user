package main

import (
	"fmt"

	"github.com/aiyi/go-user/model"
)

func add() {
	fmt.Println(model.AddEmail("email1@xxx.com", []byte("password"), []byte("salt"), 0))
	fmt.Println(model.AddPhone("18888888888", []byte("password"), []byte("salt"), 0))
	fmt.Println(model.AddQQ("openid", "nickname", 0))
	fmt.Println(model.AddWechat("openid", "nickname", 0))
	fmt.Println(model.AddWeibo("openid", "nickname", 0))
}
