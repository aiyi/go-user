package checkcode

import (
	"fmt"
)

func sendCodeToPhone(phone, checkcode string) (err error) {
	fmt.Println("sendCodeToPhone:", checkcode)
	return nil
}

func sendCodeToEmail(email, checkcode string) (err error) {
	fmt.Println("sendCodeToEmail:", checkcode)
	return nil
}
