package session

import (
	"encoding/json"
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"

	"github.com/aiyi/go-user/frontend/errors"
	"github.com/aiyi/go-user/mc"
)

const CheckCodeMaxTimes = 3 // 校验码最多可以校验多少次

type CheckCode struct {
	Key   string `json:"key"`   // email, phone...
	Code  string `json:"code"`  // 校验码
	Times int    `json:"times"` // 已经校验的次数
}

type Session struct {
	TokenSignature string `json:"token_sign"`             // Token 签名; 安全考虑, 比对客户端传过来的 token 的签名部分
	UserId         int64  `json:"user_id,omitempty"`      // 当前用户
	PasswordTag    string `json:"password_tag,omitempty"` // 认证时的 password_tag, 对于 AuthType 是 AuthTypeEmailPassword, AuthTypePhonePassword 时有效

	EmailCheckCode *CheckCode `json:"email_checkcode,omitempty"` // 邮箱校验码
	PhoneCheckCode *CheckCode `json:"phone_checkcode,omitempty"` // 手机校验码
	OAuth2State    string     `json:"oauth2_state,omitempty"`    // 微信公众号 oauth 登录的 state

	memcacheItem *memcache.Item `json:"-"`
}

// 获取 Session, 如果找不到返回 errors.ErrNotFound.
func Get(sid string) (ss *Session, err error) {
	item, err := mc.Client().Get(mc.SessionKey(sid))
	if err != nil {
		if err == memcache.ErrCacheMiss {
			err = errors.ErrNotFound
		}
		return
	}

	ss = &Session{}
	if err = json.Unmarshal(item.Value, ss); err != nil {
		return
	}
	ss.memcacheItem = item
	return
}

// 设置 sid-Session, 该 Session 是通过 Get 获取的,
// 如果在 Get 之后该 sid 对应的 Session 发生改变, 则返回错误.
func CompareAndSwap(sid string, ss *Session) (err error) {
	item := ss.memcacheItem
	if item == nil || item.Key != mc.SessionKey(sid) {
		return fmt.Errorf("参数 Session 不是通过 sid:%q 获取的", sid)
	}

	SessionBytes, err := json.Marshal(ss)
	if err != nil {
		return
	}

	item.Value = SessionBytes
	return mc.Client().CompareAndSwap(item)
}

// 添加 sid-Session, 如果 sid 冲突则返回错误.
func Add(sid string, ss *Session) (err error) {
	SessionBytes, err := json.Marshal(ss)
	if err != nil {
		return
	}

	item := memcache.Item{
		Key:   mc.SessionKey(sid),
		Value: SessionBytes,
	}
	return mc.Client().Add(&item)
}

// 无条件设置 sid-Session
func Set(sid string, ss *Session) (err error) {
	SessionBytes, err := json.Marshal(ss)
	if err != nil {
		return
	}

	item := memcache.Item{
		Key:   mc.SessionKey(sid),
		Value: SessionBytes,
	}
	return mc.Client().Set(&item)
}

// 删除 Session, 如果没有匹配则返回 errors.ErrNotFound.
func Delete(sid string) (err error) {
	if err = mc.Client().Delete(mc.SessionKey(sid)); err == memcache.ErrCacheMiss {
		err = errors.ErrNotFound
	}
	return
}
