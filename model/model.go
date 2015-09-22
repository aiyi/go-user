package model

type AuthType uint64 // 位图, 保险起见最多使用低 63 位

const (
	AuthTypeEmail  AuthType = 1 << iota // 邮箱-密码
	AuthTypePhone                       // 手机-密码, 手机-短信
	AuthTypeQQ                          // QQ
	AuthTypeWechat                      // 微信
	AuthTypeWeibo                       // 微博

	AuthTypeMask = 0x7FFFFFFFFFFFFFFF
)

var emptyByteSlice = []byte{}
