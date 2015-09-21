package model

import (
	"fmt"

	"github.com/aiyi/go-user/db"
)

// 绑定邮箱新注册账户到已经存在的账户, 密码以原账户为准.
//  调用该函数前, 请确认:
//  1. toUserId 存在并且 has_fixed
//  2. toUserId 未绑定邮箱
//  3. fromUserId 存在并且没有 has_fixed
func BindExistEmail(toUserId, fromUserId int64) (err error) {
	para := struct {
		ToUserId   int64 `sqlx:"to_user_id"`
		FromUserId int64 `sqlx:"from_user_id"`
		AuthType   int64 `sqlx:"auth_type"`
	}{
		ToUserId:   toUserId,
		FromUserId: fromUserId,
		AuthType:   AuthTypeEmail,
	}

	tx, err := db.GetDB().Beginx()
	if err != nil {
		return
	}

	// user 更新 ToUserId
	stmt1, err := tx.PrepareNamed("update user set auth_types = auth_types|:auth_type where id=:to_user_id and has_fixed=1 and auth_types&:auth_type=0")
	if err != nil {
		tx.Rollback()
		return
	}
	rslt, err := stmt1.Exec(para)
	if err != nil {
		tx.Rollback()
		return
	}
	rowsAffected, err := rslt.RowsAffected()
	if err != nil {
		tx.Rollback()
		return
	}
	if rowsAffected != 1 {
		err = fmt.Errorf("绑定用户 %d 到用户 %d 失败", para.FromUserId, para.ToUserId)
		tx.Rollback()
		return
	}

	// user 删除 FromUserId
	stmt2, err := tx.PrepareNamed("delete from user where id=:from_user_id and hax_fixed=0")
	if err != nil {
		tx.Rollback()
		return
	}
	rslt, err = stmt2.Exec(para)
	if err != nil {
		tx.Rollback()
		return
	}
	rowsAffected, err = rslt.RowsAffected()
	if err != nil {
		tx.Rollback()
		return
	}
	if rowsAffected != 1 {
		err = fmt.Errorf("绑定用户 %d 到用户 %d 失败", para.FromUserId, para.ToUserId)
		tx.Rollback()
		return
	}

	// user_email 更新 item
	stmt3, err := tx.PrepareNamed("update user_email set user_id=:to_user_id, has_fixed=1 where user_id=:from_user_id and has_fixed=0")
	if err != nil {
		tx.Rollback()
		return
	}
	rslt, err = stmt3.Exec(para)
	if err != nil {
		tx.Rollback()
		return
	}
	rowsAffected, err = rslt.RowsAffected()
	if err != nil {
		tx.Rollback()
		return
	}
	if rowsAffected != 1 {
		err = fmt.Errorf("绑定用户 %d 到用户 %d 失败", para.FromUserId, para.ToUserId)
		tx.Rollback()
		return
	}

	err = tx.Commit()
	return
}
