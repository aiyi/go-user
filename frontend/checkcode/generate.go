package checkcode

import (
	"math/rand"
	"time"
)

const digits = "0123456789"

var mathRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// 生成6位随机数字
func generateCode() string {
	checkcode := make([]byte, 6)

	timestamp := time.Now().UnixNano() / 100
	for i := 0; i < 6; i++ {
		index := int(mathRand.Int63()^timestamp) % 10
		if index < 0 {
			index = -index
		}
		checkcode[i] = digits[index]
	}

	return string(checkcode)
}
