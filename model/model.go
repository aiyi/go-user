package model

type AuthType uint64 // 位图, 保险起见最多使用低 63 位

const (
	AuthTypeEmail  AuthType = 1 << iota // 邮箱-密码
	AuthTypePhone                       // 手机-密码, 手机-短信
	AuthTypeQQ                          // QQ
	AuthTypeWechat                      // 微信
	AuthTypeWeibo                       // 微博

	AuthTypeMask = AuthType(int64(-1) ^ (int64(-1) << 63))
)

var emptyByteSlice = []byte{}
