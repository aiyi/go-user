package model

import (
	"crypto/hmac"
	"crypto/sha1"

	"github.com/chanxuehong/util/random"
)

var PasswordSalt = NewSalt() // 无主地 salt, 用于安全授权

func NewSalt() []byte {
	rd := random.NewRandom()
	return rd[:]
}

func NewPasswordTag() []byte {
	return random.NewRandomEx()
}

func EncryptPassword(password, salt []byte) []byte {
	Hash := hmac.New(sha1.New, salt)
	Hash.Write(password)
	return Hash.Sum(nil)
}
