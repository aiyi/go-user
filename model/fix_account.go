package model

import (
	"github.com/aiyi/go-user/db"
)

// 确认邮箱注册新账户
func FixEmail(userId int64) (err error) {
	stmt, err := db.GetStmt("update user as A, user_email as B set A.has_fixed=1, B.has_fixed=1 where A.id=? and A.has_fixed=0 and B.user_id=A.id and B.has_fixed=0")
	if err != nil {
		return
	}
	_, err = stmt.Exec(userId)
	return
}

// 确认手机注册新账户
func FixPhone(userId int64) (err error) {
	stmt, err := db.GetStmt("update user as A, user_phone as B set A.has_fixed=1, B.has_fixed=1 where A.id=? and A.has_fixed=0 and B.user_id=A.id and B.has_fixed=0")
	if err != nil {
		return
	}
	_, err = stmt.Exec(userId)
	return
}

// 确认QQ注册新账户
func FixQQ(userId int64) (err error) {
	stmt, err := db.GetStmt("update user as A, user_qq as B set A.has_fixed=1, B.has_fixed=1 where A.id=? and A.has_fixed=0 and B.user_id=A.id and B.has_fixed=0")
	if err != nil {
		return
	}
	_, err = stmt.Exec(userId)
	return
}

// 确认微信注册新账户
func FixWechat(userId int64) (err error) {
	stmt, err := db.GetStmt("update user as A, user_wechat as B set A.has_fixed=1, B.has_fixed=1 where A.id=? and A.has_fixed=0 and B.user_id=A.id and B.has_fixed=0")
	if err != nil {
		return
	}
	_, err = stmt.Exec(userId)
	return
}

// 确认微博注册新账户
func FixWeibo(userId int64) (err error) {
	stmt, err := db.GetStmt("update user as A, user_weibo as B set A.has_fixed=1, B.has_fixed=1 where A.id=? and A.has_fixed=0 and B.user_id=A.id and B.has_fixed=0")
	if err != nil {
		return
	}
	_, err = stmt.Exec(userId)
	return
}
