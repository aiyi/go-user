package main

import (
	"log"

	"github.com/aiyi/go-user/model"
)

func verify1() {
	//email	 757120565977810729
	//phone	 757120566451767081
	//qq	 757120566976055081
	//wechat 757120567823304489
	//weibo	 757120568624416553
	log.Println(model.VerifyEmail(757120568624416553))
	log.Println(model.VerifyPhone(757120565977810729))
	log.Println(model.VerifyQQ(757120566451767081))
	log.Println(model.VerifyWechat(757120566976055081))
	log.Println(model.VerifyWeibo(757120567823304489))
}

func verify2() {
	//email	 757120565977810729
	//phone	 757120566451767081
	//qq	 757120566976055081
	//wechat 757120567823304489
	//weibo	 757120568624416553
	log.Println(model.VerifyEmail(757120565977810729))
	log.Println(model.VerifyPhone(757120566451767081))
	log.Println(model.VerifyQQ(757120566976055081))
	log.Println(model.VerifyWechat(757120567823304489))
	log.Println(model.VerifyWeibo(757120568624416553))
}
