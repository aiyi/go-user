package main

import (
	"fmt"
	timex "github.com/chanxuehong/util/time"

	"github.com/aiyi/go-user/model"
)

func main() {
	para := model.AddByEmailParams{
		AuthType:   model.AuthTypeEmail,
		Email:      "test1@test.com",
		Password:   []byte("password"),
		Salt:       []byte("salt"),
		CreateTime: timex.Now(),
	}
	fmt.Println(model.AddByEmail(&para))
}
