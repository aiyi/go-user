package model

const (
	AuthTypeEmail  = 1 << iota // 邮箱-密码
	AuthTypePhone              // 手机-密码
	AuthTypeWechat             // 微信
	AuthTypeQQ                 // QQ
	AuthTypeWeibo              // 微博
)

var emptyByteSlice = []byte{}
