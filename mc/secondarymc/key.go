package secondarymc

// 临时 session 的 mc-key
func SessionCacheKey(sid string) string {
	return "session2/" + sid
}
