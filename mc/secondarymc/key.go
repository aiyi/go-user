package secondarymc

// 临时 token 的 mc-key
func TokenCacheKey(token string) string {
	return "token2/" + token
}
