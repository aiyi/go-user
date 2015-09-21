package main

import (
	"log"

	"github.com/aiyi/go-user/model"
)

func fix1() {
	//email	 757120565977810729
	//phone	 757120566451767081
	//qq	 757120566976055081
	//wechat 757120567823304489
	//weibo	 757120568624416553
	log.Println(model.FixEmail(757120568624416553))
	log.Println(model.FixPhone(757120565977810729))
	log.Println(model.FixQQ(757120566451767081))
	log.Println(model.FixWechat(757120566976055081))
	log.Println(model.FixWeibo(757120567823304489))
}

func fix2() {
	//email	 757120565977810729
	//phone	 757120566451767081
	//qq	 757120566976055081
	//wechat 757120567823304489
	//weibo	 757120568624416553
	log.Println(model.FixEmail(757120565977810729))
	log.Println(model.FixPhone(757120566451767081))
	log.Println(model.FixQQ(757120566976055081))
	log.Println(model.FixWechat(757120567823304489))
	log.Println(model.FixWeibo(757120568624416553))
}
