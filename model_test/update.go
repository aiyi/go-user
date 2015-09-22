package main

import (
	"log"

	"github.com/aiyi/go-user/model"
)

func update() {
	//email	 757120565977810729
	//phone	 757120566451767081
	//qq	 757120566976055081
	//wechat 757120567823304489
	//weibo	 757120568624416553
	log.Println(model.UpdateEmail(757120565977810729, "update@xxx.com"))
	log.Println(model.UpdatePhone(757120566451767081, "17777777777"))
	log.Println(model.UpdateQQ(757120566976055081, "update_qq_openid", "update_qq_nickname"))
	log.Println(model.UpdateWechat(757120567823304489, "update_wechat_openid", "update_wechat_nickname"))
	log.Println(model.UpdateWeibo(757120568624416553, "update_weibo_openid", "update_weibo_nickname"))
}
