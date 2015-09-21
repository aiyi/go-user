package main

import (
	"log"

	"github.com/aiyi/go-user/model"
)

func bind() {
	//email	 757120565977810729
	//phone	 757120566451767081
	//qq	 757120566976055081
	//wechat 757120567823304489
	//weibo	 757120568624416553
	log.Println(model.BindEmail(757120568624416553, "bind@xxx.com"))
	log.Println(model.BindPhone(757120565977810729, "19999999999"))
	log.Println(model.BindQQ(757120566451767081, "bind_qq_openid", "bind_qq_nickname"))
	log.Println(model.BindWechat(757120566976055081, "bind_wechat_openid", "bind_wechat_nickname"))
	log.Println(model.BindWeibo(757120567823304489, "bind_weibo_openid", "bind_weibo_nickname"))
}

func bind2() {
	//email	 757120565977810729
	//phone	 757120566451767081
	//qq	 757120566976055081
	//wechat 757120567823304489
	//weibo	 757120568624416553
	log.Println(model.BindEmail(757120567823304489, "xbind@xxx.com"))
	log.Println(model.BindPhone(757120568624416553, "x19999999999"))
	log.Println(model.BindQQ(757120565977810729, "xbind_qq_openid", "bind_qq_nickname"))
	log.Println(model.BindWechat(757120566451767081, "xbind_wechat_openid", "bind_wechat_nickname"))
	log.Println(model.BindWeibo(757120566976055081, "xbind_weibo_openid", "bind_weibo_nickname"))
}
