package model

const (
	AuthTypeEmail         = 1 << iota // 邮箱-密码
	AuthTypePhonePassword             // 手机-密码
	AuthTypePhoneCode                 // 手机-短信验证码
	AuthTypeWechat                    // 微信
	AuthTypeQQ                        // QQ
	AuthTypeWeibo                     // 微博
)

var emptyByteSlice = []byte{}
