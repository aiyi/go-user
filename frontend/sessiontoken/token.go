package sessiontoken

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/chanxuehong/util/security"
)

const (
	AuthTypeEmailPassword = "email_password" // 邮箱+密码
	AuthTypeEmailCaptcha  = "email_captcha"  // 邮箱+验证码, 验证码推送到邮箱
	AuthTypePhonePassword = "phone_password" // 手机+密码
	AuthTypePhoneCaptcha  = "phone_captcha"  // 手机+验证码, 验证码短信推送给手机
	AuthTypeOAuthQQ       = "oauth_qq"       // QQ oauth
	AuthTypeOAuthWechat   = "oauth_wechat"   // 微信 oauth
	AuthTypeOAuthWeibo    = "oauth_weibo"    // 微博 oauth
)

// 该结构体序列化后, 在客户端和服务器之间传递, 作为认证标识
type SessionToken struct {
	SessionId string `json:"sid"`

	// 以下为 token 相关部分
	Id            string `json:"id"`            // token 的标识, 每次刷新 token 改变此值
	Authenticated bool   `json:"authenticated"` // 是否是认证后的 token, 否则为临时 token
	UserId        int64  `json:"user_id"`       // token 的拥有者
	AuthType      string `json:"auth_type"`     // token 的认证类型
	PasswordTag   string `json:"password_tag"`  // 认证时的 password_tag, 对于 AuthType 是 AuthTypeEmailPassword, AuthTypePhonePassword 时有效
	ExpAccess     int64  `json:"exp_access"`    // 该 token 的过期时间
	ExpRefresh    int64  `json:"exp_refresh"`   // 通过该 token 换取新的 token 的截至时间, 固定值, 不会变化
}

func (token *SessionToken) Encode(securityKey []byte) ([]byte, error) {
	jsonBytes, err := json.Marshal(token)
	if err != nil {
		return nil, err
	}
	Hash := hmac.New(sha1.New, securityKey)

	n1 := base64.URLEncoding.EncodedLen(len(jsonBytes))
	n2 := hex.EncodedLen(Hash.Size())
	buf := make([]byte, n1+1+n2)

	base64.URLEncoding.Encode(buf, jsonBytes)
	buf[n1] = '|'

	Hash.Write(buf[:n1])
	hex.Encode(buf[n1+1:], Hash.Sum(nil))

	return buf, nil
}

var tokenBytesSplitSep = []byte{'|'}

func (token *SessionToken) Decode(tokenBytes []byte, securityKey []byte) error {
	bytesArr := bytes.Split(tokenBytes, tokenBytesSplitSep)
	if len(bytesArr) != 2 {
		return errors.New("invalid token input")
	}

	Hash := hmac.New(sha1.New, securityKey)
	Hash.Write(bytesArr[0])

	HashSum := make([]byte, hex.EncodedLen(Hash.Size()))
	hex.Encode(HashSum, Hash.Sum(nil))

	buf := make([]byte, base64.URLEncoding.DecodedLen(len(bytesArr[0])))

	if !security.SecureCompare(HashSum, bytesArr[1]) {
		n, err := base64.URLEncoding.Decode(buf, bytesArr[0])
		if err != nil {
			return errors.New("invalid token input")
		}
		if err = json.Unmarshal(buf[:n], token); err != nil {
			return errors.New("invalid token input")
		}
		return errors.New("invalid token input")
	}
	n, err := base64.URLEncoding.Decode(buf, bytesArr[0])
	if err != nil {
		return err
	}
	if err = json.Unmarshal(buf[:n], token); err != nil {
		return err
	}
	return nil
}
