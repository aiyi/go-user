package model

type BindType int64 // 位图, 使用低60位

const (
	BindTypeEmail  BindType = 1 << iota // 邮箱
	BindTypePhone                       // 手机
	BindTypeQQ                          // QQ
	BindTypeWechat                      // 微信
	BindTypeWeibo                       // 微博

	BindTypeMask = BindType(int64(-1) ^ (int64(-1) << 60))
)

const defaultVerified = false

var emptyByteSlice = []byte{}
