package checkcode

import (
	"errors"
	"math/rand"
	"time"
)

const (
	checkcodeLength = 6

	digits = "0123456789"
)

var mathRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func newCheckcode() []byte {
	timestamp := time.Now().UnixNano() / 100

	checkcode := make([]byte, checkcodeLength)
	for i := 0; i < checkcodeLength; i++ {
		index := int(mathRand.Int63()^timestamp) % 10
		if index < 0 {
			index = -index
		}
		checkcode[i] = digits[index]
	}

	return checkcode
}

func sendToPhone(phone []byte, checkcode string) (err error) {
	return errors.New("not supported")
}

func sendToEmail(email []byte, checkcode string) (err error) {
	return errors.New("not supported")
}
