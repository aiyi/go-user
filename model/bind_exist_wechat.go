package model

import (
	"errors"
	"fmt"

	"github.com/aiyi/go-user/db"
)

// 绑定微信新注册账户到已经存在的账户, 密码以原账户为准.
//  调用该函数前, 请确认:
//  1. toUserId != fromUserId
//  2. toUserId 存在并且 verified
//  3. fromUserId 存在并且没有 verified
//  4. toUserId 未绑定微信
//  5. fromUserId 是微信新注册账户
func BindExistWechat(toUserId, fromUserId int64) (err error) {
	if toUserId == fromUserId {
		return errors.New("toUserId 不能等于 fromUserId")
	}

	para := struct {
		ToUserId   int64    `sqlx:"to_user_id"`
		FromUserId int64    `sqlx:"from_user_id"`
		AuthType   AuthType `sqlx:"auth_type"`
	}{
		ToUserId:   toUserId,
		FromUserId: fromUserId,
		AuthType:   AuthTypeWechat,
	}

	tx, err := db.GetDB().Beginx()
	if err != nil {
		return
	}

	// user 更新 ToUserId
	stmt1, err := tx.PrepareNamed("update user set auth_types = auth_types|:auth_type where id=:to_user_id and verified=1 and auth_types&:auth_type=0")
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
	stmt2, err := tx.PrepareNamed("delete from user where id=:from_user_id and verified=0 and auth_types=:auth_type")
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

	// user_wechat 更新 item
	stmt3, err := tx.PrepareNamed("update user_wechat set user_id=:to_user_id, verified=1 where user_id=:from_user_id and verified=0")
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
