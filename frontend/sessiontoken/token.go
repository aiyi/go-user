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
	AuthTypeGuest          = "guest"           // 游客
	AuthTypeEmailPassword  = "email_password"  // 邮箱+密码
	AuthTypeEmailCheckcode = "email_checkcode" // 邮箱+校验码, 校验码推送到邮箱
	AuthTypePhonePassword  = "phone_password"  // 手机+密码
	AuthTypePhoneCheckcode = "phone_checkcode" // 手机+校验码, 校验码短信推送给手机
	AuthTypeOAuthQQ        = "oauth_qq"        // QQ oauth
	AuthTypeOAuthWechat    = "oauth_wechat"    // 微信 oauth
	AuthTypeOAuthWeibo     = "oauth_weibo"     // 微博 oauth
)

func ExpirationAccess(timestamp int64) int64 {
	return timestamp + 7200
}

func ExpirationRefresh(timestamp int64) int64 {
	return timestamp + 31556952
}

// 客户端访问 API 的令牌, 客户端和服务器交互的数据结构
type SessionToken struct {
	SessionId         string `json:"sid"`         // 服务器索引 Session 的 key
	TokenId           string `json:"token_id"`    // token 的标识, 每次刷新 token 改变此值
	AuthType          string `json:"auth_type"`   // token 的认证类型
	ExpirationAccess  int64  `json:"exp_access"`  // 该 token 的过期时间
	ExpirationRefresh int64  `json:"exp_refresh"` // 刷新 token 的截至时间, 固定值, 不会变化

	Signatrue string `json:"-"` // 和客户端交互的 token 签名部分; 在 SessionToken.Encode 或者 SessionToken.Decode 才会获取到正确的值
}

// trim(url_base64(json(token))) + "." + hex(hmac-sha256(base64_str))
func (token *SessionToken) Encode(securityKey []byte) ([]byte, error) {
	const signatureLen = 64 // hmac-sha256

	jsonBytes, err := json.Marshal(token)
	if err != nil {
		return nil, err
	}

	base64BytesLen := base64.URLEncoding.EncodedLen(len(jsonBytes))
	buf := make([]byte, base64BytesLen+1+signatureLen)

	base64Bytes := buf[:base64BytesLen]
	base64.URLEncoding.Encode(base64Bytes, jsonBytes)

	// 去掉 base64 编码尾部的 '='
	base64Bytes = base64Trim(base64Bytes)
	base64BytesLen = len(base64Bytes)
	signatureOffset := base64BytesLen + 1
	buf = buf[:signatureOffset+signatureLen]

	buf[base64BytesLen] = '.'

	signature := buf[signatureOffset:]
	Hash := hmac.New(sha256.New, securityKey)
	Hash.Write(base64Bytes)
	hex.Encode(signature, Hash.Sum(nil))
	token.Signatrue = string(signature)

	return buf, nil
}

var tokenBytesSplitSep = []byte{'.'}

// trim(url_base64(json(token))) + "." + hex(hmac-sha256(base64_str))
func (token *SessionToken) Decode(tokenBytes []byte, securityKey []byte) error {
	const signatureLen = 64 // hmac-sha256

	bytesArray := bytes.Split(tokenBytes, tokenBytesSplitSep)
	if len(bytesArray) < 2 {
		return errors.New("invalid token bytes")
	}

	// 验证签名
	signatrue := make([]byte, signatureLen)
	Hash := hmac.New(sha256.New, securityKey)
	Hash.Write(bytesArray[0])
	hex.Encode(signatrue, Hash.Sum(nil))
	if !bytes.Equal(signatrue, bytesArray[1]) {
		return errors.New("invalid token bytes, signature mismatch")
	}

	// 解码
	temp := signatrue[:4]                       // signatrue 不再使用, 利用其空间
	copy(temp, tokenBytes[len(bytesArray[0]):]) // 保护 tokenBytes
	defer func() {
		copy(tokenBytes[len(bytesArray[0]):], temp) // 恢复 tokenBytes
		token.Signatrue = string(bytesArray[1])
	}()

	base64Bytes := base64Pad(bytesArray[0])
	base64Decoder := base64.NewDecoder(base64.URLEncoding, bytes.NewReader(base64Bytes))
	return json.NewDecoder(base64Decoder).Decode(token)
}
