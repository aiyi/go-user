package model

import (
	"errors"
	"fmt"

	"github.com/aiyi/go-user/db"
)

// 绑定微博新注册账户到已经存在的账户, 密码以原账户为准.
//  调用该函数前, 请确认:
//  1. toUserId 存在并且 has_fixed
//  2. fromUserId 存在并且没有 has_fixed
//  3. toUserId 未绑定微博
func BindExistWeibo(toUserId, fromUserId int64) (err error) {
	if toUserId == fromUserId {
		return errors.New("toUserId 不能等于 fromUserId")
	}

	para := struct {
		ToUserId   int64 `sqlx:"to_user_id"`
		FromUserId int64 `sqlx:"from_user_id"`
		AuthType   int64 `sqlx:"auth_type"`
	}{
		ToUserId:   toUserId,
		FromUserId: fromUserId,
		AuthType:   AuthTypeWeibo,
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
	rslt1, err := stmt1.Exec(para)
	if err != nil {
		tx.Rollback()
		return
	}
	rowsAffected1, err := rslt1.RowsAffected()
	if err != nil {
		tx.Rollback()
		return
	}

	// user 删除 FromUserId
	stmt2, err := tx.PrepareNamed("delete from user where id=:from_user_id and hax_fixed=0")
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

	// user_weibo 更新 item
	stmt3, err := tx.PrepareNamed("update user_weibo set user_id=:to_user_id, has_fixed=1 where user_id=:from_user_id and has_fixed=0")
	if err != nil {
		tx.Rollback()
		return
	}
	rslt3, err := stmt3.Exec(para)
	if err != nil {
		tx.Rollback()
		return
	}
	rowsAffected3, err := rslt3.RowsAffected()
	if err != nil {
		tx.Rollback()
		return
	}

	if rowsAffected1 != rowsAffected2 || rowsAffected1 != rowsAffected3 {
		err = fmt.Errorf("绑定用户 %d 到用户 %d 失败", para.FromUserId, para.ToUserId)
		tx.Rollback()
		return
	}

	err = tx.Commit()
	return
}
