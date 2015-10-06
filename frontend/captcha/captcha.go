package captcha

import (
	"errors"
	"math/rand"
	"time"
)

const (
	captchaLength = 6

	digits = "0123456789"
)

var mathRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func newCaptcha() []byte {
	timestamp := time.Now().UnixNano() / 100

	captcha := make([]byte, captchaLength)
	for i := 0; i < captchaLength; i++ {
		index := int(mathRand.Int63()^timestamp) % 10
		if index < 0 {
			index = -index
		}
		captcha[i] = digits[index]
	}

	return captcha
}

func sendToPhone(phone []byte, captcha string) (err error) {
	return errors.New("not supported")
}

func sendToEmail(email []byte, captcha string) (err error) {
	return errors.New("not supported")
}
