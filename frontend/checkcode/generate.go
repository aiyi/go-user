package checkcode

import (
	"math/rand"
	"time"
)

var mathRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// 生成6位随机数字
func generateCode() string {
	const (
		digits  = "0123456789"
		digits2 = "123456789"
	)
	timestamp := time.Now().Nanosecond() / 100
	checkcode := make([]byte, 6)
	checkcode[0] = digits2[(mathRand.Int()^timestamp)%len(digits2)]
	for i := 1; i < 6; i++ {
		checkcode[i] = digits[(mathRand.Int()^timestamp)%len(digits)]
	}
	return string(checkcode)
}
