package model

import (
	"fmt"

	"github.com/aiyi/go-user/db"
)

// 确认微信注册新账户
func FixWechat(userId int64) (err error) {
	tx, err := db.GetDB().Begin()
	if err != nil {
		return
	}

	// user_wechat 表更新 item
	stmt1, err := tx.Prepare("update user_wechat set has_fixed=1 where user_id=? and has_fixed=0")
	if err != nil {
		tx.Rollback()
		return
	}
	rslt1, err := stmt1.Exec(userId)
	if err != nil {
		tx.Rollback()
		return
	}
	rowsAffected1, err := rslt1.RowsAffected()
	if err != nil {
		tx.Rollback()
		return
	}

	// user 表更新 item
	stmt2, err := tx.Prepare("update user set has_fixed=1 where id=? and has_fixed=0")
	if err != nil {
		tx.Rollback()
		return
	}
	rslt2, err := stmt2.Exec(userId)
	if err != nil {
		tx.Rollback()
		return
	}
	rowsAffected2, err := rslt2.RowsAffected()
	if err != nil {
		tx.Rollback()
		return
	}

	if rowsAffected1 != rowsAffected2 {
		err = fmt.Errorf("确认新建账户失败, ID为 %d", userId)
		tx.Rollback()
		return
	}

	err = tx.Commit()
	return
}
