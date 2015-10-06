package authcode

import (
	"math/rand"
	"time"
)

const (
	CodeLength = 6

	digits = "0123456789"
)

var mathRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// 生成6位随机数字校验码
func NewCode() []byte {
	timestamp := time.Now().UnixNano() / 100
	bs := make([]byte, CodeLength)
	for i := 0; i < CodeLength; i++ {
		index := int(mathRand.Int63()^timestamp) % 10
		if index < 0 {
			index = -index
		}
		bs[i] = digits[index]
	}
	return bs
}
