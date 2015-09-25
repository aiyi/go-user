package mc

import (
	"strconv"
)

func UserCacheKey(userId int64) string {
	return "user/" + strconv.FormatInt(userId, 16)
}
