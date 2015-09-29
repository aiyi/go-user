package frontend

import (
	"github.com/aiyi/go-user/model"
)

type AuthType int64 // 认证类型, 低 60 位对应绑定类型, 见 model.BindType

func (typ AuthType) BindType() model.BindType {
	return model.BindType(typ) & model.BindTypeMask
}

const (
	AuthTypeEmailPassword = AuthType(model.BindTypeEmail)  // 邮箱-密码
	AuthTypePhonePassword = AuthType(model.BindTypePhone)  // 手机-密码
	AuthTypeQQ            = AuthType(model.BindTypeQQ)     // QQ oauth2
	AuthTypeWechat        = AuthType(model.BindTypeWechat) // 微信 oauth2
	AuthTypeWeibo         = AuthType(model.BindTypeWeibo)  // 微博 oauth2

	AuthTypeEmailCaptcha = AuthTypeEmailPassword | 1<<60 // 邮箱-验证码
	AuthTypePhoneCaptcha = AuthTypePhonePassword | 1<<60 // 手机-验证码
)

type Token struct {
	Value         string   `json:"value"`         // token 的值
	Authenticated bool     `json:"authenticated"` // 该 token 是否认证过, 没有认证的就是临时 token, 没有实际业务权限, 保存状态而已
	UserId        int64    `json:"user_id"`       // token 的拥有者
	AuthType      AuthType `json:"auth_type"`     // token 的认证类型
	PasswordTag   string   `json:"password_tag"`  // 认证时的 password_tag, 对于 AuthType 是 AuthTypeEmailPassword, AuthTypePhonePassword 时有效
	ExpAccess     int64    `json:"exp_access"`    // 该 token 的过期时间
	ExpRefresh    int64    `json:"exp_refresh"`   // 通过该 token 换取新的 token 的截至时间, 固定值, 不会变化
}
