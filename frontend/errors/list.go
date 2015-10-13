package errors

const (
	ErrCodeOK = 0

	ErrCodeBadRequest          = 10000 // 输入参数不正确, 包括 url 参数和 http-body-json 参数
	ErrCodeInternalServerError = 10001 // 内部服务器出错

	ErrCodeTokenMissing        = 20000 // token 缺失
	ErrCodeTokenEncodeFailed   = 20001 // token 编码出错
	ErrCodeTokenDecodeFailed   = 20002 // token 解码出错
	ErrCodeTokenAccessExpired  = 20003 // token access 过期
	ErrCodeTokenRefreshExpired = 20004 // token refresh 过期
	ErrCodeTokenShouldNotGuest = 20005 // token 不能是 guest 认证的
	ErrCodeTokenInvalid        = 20006 // token 失效, 应该重新认证

	ErrCodeAuthFailed         = 20010 // 认证失败
	ErrCodeOAuthStateMismatch = 20011 // oauth state 不匹配
)

var (
	ErrOK = &Error{
		ErrCode: ErrCodeOK,
		ErrMsg:  "success",
	}

	ErrBadRequest = &Error{
		ErrCode: ErrCodeBadRequest,
		ErrMsg:  "bad request",
	}
	ErrInternalServerError = &Error{
		ErrCode: ErrCodeInternalServerError,
		ErrMsg:  "internal server error",
	}

	ErrTokenMissing = &Error{
		ErrCode: ErrCodeTokenMissing,
		ErrMsg:  "token missing",
	}
	ErrTokenEncodeFailed = &Error{
		ErrCode: ErrCodeTokenEncodeFailed,
		ErrMsg:  "token encoding failed",
	}
	ErrTokenDecodeFailed = &Error{
		ErrCode: ErrCodeTokenDecodeFailed,
		ErrMsg:  "token decoding failed",
	}
	ErrTokenAccessExpired = &Error{
		ErrCode: ErrCodeTokenAccessExpired,
		ErrMsg:  "token access expired",
	}
	ErrTokenRefreshExpired = &Error{
		ErrCode: ErrCodeTokenRefreshExpired,
		ErrMsg:  "token refresh expired",
	}
	ErrTokenShouldNotGuest = &Error{
		ErrCode: ErrCodeTokenShouldNotGuest,
		ErrMsg:  "token should not be authenticated via guest",
	}
	ErrTokenInvalid = &Error{
		ErrCode: ErrCodeTokenInvalid,
		ErrMsg:  "token is invalid",
	}

	ErrAuthFailed = &Error{
		ErrCode: ErrCodeAuthFailed,
		ErrMsg:  "authentication failed",
	}
	ErrOAuthStateMismatch = &Error{
		ErrCode: ErrCodeOAuthStateMismatch,
		ErrMsg:  "the oauth state mismatch",
	}
)
