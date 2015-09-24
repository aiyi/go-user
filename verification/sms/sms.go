package sms

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/chanxuehong/util/security"

	"github.com/aiyi/go-user/mc"
	"github.com/aiyi/go-user/verification"
)

var compareByteSlice = make([]byte, verification.CodeLength) // 安全比较需要

func key(phone string) string {
	return "phone:" + phone
}

// 生成 phone 对应的 code, 6位随机数字, 有效期为 300 秒.
func NewCode(phone string) ([]byte, error) {
	item := memcache.Item{
		Key:        key(phone),
		Value:      verification.NewCode(),
		Expiration: 300,
	}
	if err := mc.Client().Set(&item); err != nil {
		return nil, err
	}
	return item.Value, nil
}

// 验证 code
func VerifyCode(phone string, code []byte) (bool, error) {
	item, err := mc.Client().Get(key(phone))
	if err != nil && err != memcache.ErrCacheMiss { // 出错, 直接返回错误
		return false, err
	}
	if err == memcache.ErrCacheMiss {
		if !security.SecureCompare(compareByteSlice, compareByteSlice) {
			return false, nil
		}
		return false, nil
	}
	if !security.SecureCompare(code, item.Value) {
		return false, nil
	}
	return true, nil
}
