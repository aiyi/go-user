package token

import (
	"encoding/json"
	"errors"

	"github.com/bradfitz/gomemcache/memcache"

	"github.com/aiyi/go-user/frontend"
	"github.com/aiyi/go-user/mc"
	"github.com/aiyi/go-user/mc/secondarymc"
)

type Session struct {
	SessionToken SessionToken `json:"session_token"` // SessionToken 副本
	EmailCaptcha string       `json:"email_captcha"` // 邮箱验证码
	PhoneCaptcha string       `json:"phone_captcha"` // 短信验证码
}

// 获取 Session, 如果找不到返回 frontend.ErrNotFound.
func SessionGet(token *SessionToken) (*Session, error) {
	if token == nil {
		return nil, errors.New("nil SessionToken")
	}

	var (
		memcacheClient  *memcache.Client
		memcacheItemKey string
	)
	if token.AuthType == AuthTypeGuest {
		memcacheClient = mc.Client()
		memcacheItemKey = mc.SessionCacheKey(token.SessionId)
	} else {
		memcacheClient = secondarymc.Client()
		memcacheItemKey = secondarymc.SessionCacheKey(token.SessionId)
	}

	item, err := memcacheClient.Get(memcacheItemKey)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			err = frontend.ErrNotFound
		}
		return nil, err
	}

	var ss Session
	if err = json.Unmarshal(item.Value, &ss); err != nil {
		return nil, err
	}
	return &ss, nil
}

func SessionAdd(ss *Session) (err error) {
	if ss == nil {
		return errors.New("nil Session")
	}

	SessionBytes, err := json.Marshal(ss)
	if err != nil {
		return
	}

	if ss.SessionToken.AuthType == AuthTypeGuest {
		item := memcache.Item{
			Key:   mc.SessionCacheKey(ss.SessionToken.SessionId),
			Value: SessionBytes,
		}
		return mc.Client().Add(&item)
	} else {
		item := memcache.Item{
			Key:   secondarymc.SessionCacheKey(ss.SessionToken.SessionId),
			Value: SessionBytes,
		}
		return secondarymc.Client().Add(&item)
	}
}

func SessionSet(ss *Session) (err error) {
	if ss == nil {
		return errors.New("nil Session")
	}

	SessionBytes, err := json.Marshal(ss)
	if err != nil {
		return
	}

	if ss.SessionToken.AuthType == AuthTypeGuest {
		item := memcache.Item{
			Key:   mc.SessionCacheKey(ss.SessionToken.SessionId),
			Value: SessionBytes,
		}
		return mc.Client().Set(&item)
	} else {
		item := memcache.Item{
			Key:   secondarymc.SessionCacheKey(ss.SessionToken.SessionId),
			Value: SessionBytes,
		}
		return secondarymc.Client().Set(&item)
	}
}

// 删除 Session, 如果没有匹配则返回 frontend.ErrNotFound.
func SessionDelete(token *SessionToken) (err error) {
	if token == nil {
		return errors.New("nil SessionToken")
	}

	if token.AuthType == AuthTypeGuest {
		err = mc.Client().Delete(mc.SessionCacheKey(token.SessionId))
	} else {
		err = secondarymc.Client().Delete(secondarymc.SessionCacheKey(token.SessionId))
	}
	if err == memcache.ErrCacheMiss {
		err = frontend.ErrNotFound
	}
	return
}
