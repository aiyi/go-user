package model

const (
	AuthTypeEmail  = 1 << iota // 邮箱-密码
	AuthTypePhone              // 手机-密码, 手机-短信
	AuthTypeQQ                 // QQ
	AuthTypeWechat             // 微信
	AuthTypeWeibo              // 微博
)

var emptyByteSlice = []byte{}
