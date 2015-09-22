package main

import (
	"log"

	"github.com/aiyi/go-user/model"
)

func unbind() {
	//email	 757120565977810729
	//phone	 757120566451767081
	//qq	 757120566976055081
	//wechat 757120567823304489
	//weibo	 757120568624416553

	//	log.Println(model.UnbindEmail(757120565977810729))
	//	log.Println(model.UnbindPhone(757120565977810729))
	//	log.Println(model.UnbindQQ(757120565977810729))
	//	log.Println(model.BindEmail(757120565977810729, "unbind@test.com"))
	//	log.Println(model.UnbindQQ(757120565977810729))

	//	log.Println(model.UnbindQQ(757120566451767081))
	//	log.Println(model.UnbindWechat(757120566451767081))

	//	log.Println(model.UnbindWechat(757120566976055081))
	//	log.Println(model.UnbindWeibo(757120566976055081))

	//	log.Println(model.UnbindEmail(757120567823304489))
	//	log.Println(model.UnbindWeibo(757120567823304489))

	//	log.Println(model.UnbindEmail(757120568624416553))
	//	log.Println(model.UnbindPhone(757120568624416553))

	log.Println(model.UnbindEmail(757120565977810729))
	log.Println(model.UnbindPhone(757120566451767081))
	log.Println(model.UnbindQQ(757120566976055081))
	log.Println(model.UnbindWechat(757120567823304489))
	log.Println(model.UnbindWeibo(757120568624416553))
}
