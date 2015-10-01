package sessiontoken

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
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
	TokenId       string `json:"token_id"`      // token 的标识, 每次刷新 token 改变此值
	Authenticated bool   `json:"authenticated"` // 是否是认证后的 token, 否则为临时 token
	UserId        int64  `json:"user_id"`       // token 的拥有者
	AuthType      string `json:"auth_type"`     // token 的认证类型
	PasswordTag   string `json:"password_tag"`  // 认证时的 password_tag, 对于 AuthType 是 AuthTypeEmailPassword, AuthTypePhonePassword 时有效
	ExpAccess     int64  `json:"exp_access"`    // 该 token 的过期时间
	ExpRefresh    int64  `json:"exp_refresh"`   // 通过该 token 换取新的 token 的截至时间, 固定值, 不会变化
}

// url_base64(json(token)) + "." + hex(sign(base64_str))
func (token *SessionToken) Encode(securityKey []byte) ([]byte, error) {
	jsonBytes, err := json.Marshal(token)
	if err != nil {
		return nil, err
	}
	Hash := hmac.New(sha256.New, securityKey)

	n1 := base64.URLEncoding.EncodedLen(len(jsonBytes))
	n2 := hex.EncodedLen(Hash.Size())
	buf := make([]byte, n1+1+n2)

	base64Bytes := buf[:n1]
	base64.URLEncoding.Encode(base64Bytes, jsonBytes)
	base64Bytes = base64Trim(base64Bytes)

	n1 = len(base64Bytes)
	buf[n1] = '.'

	Hash.Write(base64Bytes)
	hex.Encode(buf[n1+1:], Hash.Sum(nil))

	return buf[:n1+1+n2], nil
}

var tokenBytesSplitSep = []byte{'.'}

// url_base64(json(token)) + "." + hex(sign(base64_str))
func (token *SessionToken) Decode(tokenBytes []byte, securityKey []byte) error {
	bytesArray := bytes.Split(tokenBytes, tokenBytesSplitSep)
	if len(bytesArray) < 2 {
		return errors.New("invalid token input")
	}

	Hash := hmac.New(sha256.New, securityKey)
	HashSum := make([]byte, hex.EncodedLen(Hash.Size()))
	Hash.Write(bytesArray[0])
	hex.Encode(HashSum, Hash.Sum(nil))

	if !bytes.Equal(HashSum, bytesArray[1]) {
		return errors.New("invalid token input, signature mismatch")
	}

	bytesArray[0] = base64Pad(bytesArray[0])
	buf := make([]byte, base64.URLEncoding.DecodedLen(len(bytesArray[0])))
	n, err := base64.URLEncoding.Decode(buf, bytesArray[0])
	if err != nil {
		return err
	}

	return json.Unmarshal(buf[:n], token)
}
