package sessiontoken

import (
	"encoding/json"

	"github.com/bradfitz/gomemcache/memcache"

	"github.com/aiyi/go-user/frontend/errors"
	"github.com/aiyi/go-user/mc"
)

type CheckCode struct {
	Code       string `json:"code"`
	Expiration int64  `json:"exp"`
}

type Session struct {
	SessionTokenSignature string     `json:"session_token_sign"`        // SessionToken 签名; 安全考虑, 比对客户端传过来的 sessiontoken 的签名部分
	UserId                int64      `json:"user_id,omitempty"`         // token 的拥有者
	PasswordTag           string     `json:"password_tag,omitempty"`    // 认证时的 password_tag, 对于 AuthType 是 AuthTypeEmailPassword, AuthTypePhonePassword 时有效
	EmailCheckCode        *CheckCode `json:"email_checkcode,omitempty"` // 邮箱校验码
	PhoneCheckCode        *CheckCode `json:"phone_checkcode,omitempty"` // 短信校验码
}

// 获取 Session, 如果找不到返回 errors.ErrNotFound.
func SessionGet(sid string) (*Session, error) {
	item, err := mc.Client().Get(mc.SessionCacheKey(sid))
	if err != nil {
		if err == memcache.ErrCacheMiss {
			err = errors.ErrNotFound
		}
		return nil, err
	}

	var ss Session
	if err = json.Unmarshal(item.Value, &ss); err != nil {
		return nil, err
	}
	return &ss, nil
}

func SessionAdd(sid string, ss *Session) (err error) {
	SessionBytes, err := json.Marshal(ss)
	if err != nil {
		return
	}

	item := memcache.Item{
		Key:   mc.SessionCacheKey(sid),
		Value: SessionBytes,
	}
	return mc.Client().Add(&item)
}

func SessionSet(sid string, ss *Session) (err error) {
	SessionBytes, err := json.Marshal(ss)
	if err != nil {
		return
	}

	item := memcache.Item{
		Key:   mc.SessionCacheKey(sid),
		Value: SessionBytes,
	}
	return mc.Client().Set(&item)
}

// 删除 Session, 如果没有匹配则返回 errors.ErrNotFound.
func SessionDelete(sid string) (err error) {
	if err = mc.Client().Delete(mc.SessionCacheKey(sid)); err == memcache.ErrCacheMiss {
		err = errors.ErrNotFound
	}
	return
}
