package checkcode

import (
	"errors"
)

func sendToPhone(phone string, checkcode []byte) (err error) {
	return errors.New("not supported")
}

func sendToEmail(email string, checkcode []byte) (err error) {
	return errors.New("not supported")
}
