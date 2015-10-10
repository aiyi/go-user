package checkcode

import (
	"errors"
)

func sendToPhone(phone string, checkcode []byte) (err error) {
	//return errors.New("not supported")
	return nil
}

func sendToEmail(email string, checkcode []byte) (err error) {
	return errors.New("not supported")
}
