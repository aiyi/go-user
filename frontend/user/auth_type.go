package user

import (
	"github.com/aiyi/go-user/frontend/token"
)

const (
	AuthTypeGuest          = token.AuthTypeGuest
	AuthTypeEmailPassword  = token.AuthTypeEmailPassword
	AuthTypeEmailCheckCode = token.AuthTypeEmailCheckCode
	AuthTypePhonePassword  = token.AuthTypePhonePassword
	AuthTypePhoneCheckCode = token.AuthTypePhoneCheckCode
	AuthTypeOAuthQQ        = token.AuthTypeOAuthQQ
	AuthTypeOAuthWechat    = token.AuthTypeOAuthWechat
	AuthTypeOAuthWeibo     = token.AuthTypeOAuthWeibo
)
