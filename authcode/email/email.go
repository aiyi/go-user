package email

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/chanxuehong/util/security"

	"github.com/aiyi/go-user/authcode"
	"github.com/aiyi/go-user/mc"
)

func key(email string) string {
	return "email:" + email
}

// 生成 email 对应的 code, 6位随机数字, 有效期为 300 秒.
func NewCode(email string) ([]byte, error) {
	item := memcache.Item{
		Key:        key(email),
		Value:      authcode.NewCode(),
		Expiration: 300,
	}
	if err := mc.Client().Set(&item); err != nil {
		return nil, err
	}
	return item.Value, nil
}

var compareByteSlice = make([]byte, authcode.CodeLength) // 安全比较需要

// 验证 code
func VerifyCode(email string, code []byte) (bool, error) {
	item, err := mc.Client().Get(key(email))
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
