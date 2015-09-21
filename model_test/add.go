package main

import (
	"fmt"

	"github.com/aiyi/go-user/model"
)

func add() {
	fmt.Println(model.AddEmail("email1@xxx.com", []byte("password"), []byte("salt")))
	fmt.Println(model.AddPhone("18888888888", []byte("password"), []byte("salt")))
	fmt.Println(model.AddQQ("openid", "nickname"))
	fmt.Println(model.AddWechat("openid", "nickname"))
	fmt.Println(model.AddWeibo("openid", "nickname"))
}
