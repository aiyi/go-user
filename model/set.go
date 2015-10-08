package model

import (
	"github.com/aiyi/go-user/db"
)

// 修改昵称
func SetNickname(userId int64, nickname string) (err error) {
	if err = removeFromCache(userId); err != nil {
		return
	}
	if err = setNickname(userId, nickname); err != nil {
		return
	}
	return syncToCache(userId)
}

func setNickname(userId int64, nickname string) (err error) {
	stmt, err := db.GetStmt("update user set nickname=? where id=? and verified=1")
	if err != nil {
		return
	}

	_, err = stmt.Exec(nickname, userId)
	return
}

// 修改密码
func SetPassword(userId int64, password, salt []byte) (err error) {
	if err = removeFromCache(userId); err != nil {
		return
	}
	if err = setPassword(userId, password, salt); err != nil {
		return
	}
	return syncToCache(userId)
}

func setPassword(userId int64, password, salt []byte) (err error) {
	stmt, err := db.GetStmt("update user set password=?, password_tag=?, salt=? where id=? and verified=1")
	if err != nil {
		return
	}

	_, err = stmt.Exec(password, NewPasswordTag(), salt, userId)
	return
}
