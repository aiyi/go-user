package model

import (
	"fmt"

	"github.com/aiyi/go-user/db"
)

// 解绑手机认证.
//  调用该函数前, 请确认:
//  1. 该用户存在并且 has_fixed
//  2. 该用户除了手机认证还有别的认证
func UnbindPhone(userId int64) (err error) {
	para := struct {
		UserId      int64    `sqlx:"user_id"`
		NotAuthType AuthType `sqlx:"not_auth_type"`
	}{
		UserId:      userId,
		NotAuthType: ^AuthTypePhone,
	}

	tx, err := db.GetDB().Beginx()
	if err != nil {
		return
	}

	// user_phone 表删除一个 item
	stmt1, err := tx.Prepare("delete from user_phone where user_id=? and has_fixed=1")
	if err != nil {
		tx.Rollback()
		return
	}
	rslt1, err := stmt1.Exec(para.UserId)
	if err != nil {
		tx.Rollback()
		return
	}
	rowsAffected1, err := rslt1.RowsAffected()
	if err != nil {
		tx.Rollback()
		return
	}

	// user 更新 item
	stmt2, err := tx.PrepareNamed("update user set auth_types = auth_types&:not_auth_type where id=:user_id and has_fixed=1 and auth_types&:not_auth_type<>0")
	if err != nil {
		tx.Rollback()
		return
	}
	rslt2, err := stmt2.Exec(para)
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
		err = fmt.Errorf("用户 %d 解绑手机失败", para.UserId)
		tx.Rollback()
		return
	}

	err = tx.Commit()
	return
}
