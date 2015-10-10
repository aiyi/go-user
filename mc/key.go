package mc

import (
	"strconv"
)

func UserKey(userId int64) string {
	return "user/" + strconv.FormatInt(userId, 16)
}

func SessionKey(sid string) string {
	return "session/" + sid
}
