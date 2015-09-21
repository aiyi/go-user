package model

import (
	"fmt"

	"github.com/aiyi/go-user/db"
)

type BindEmailParams struct {
	UserId int64  `sqlx:"user_id"` // 绑定到这个用户
	Email  string `sqlx:"email"`
}

// 给用户绑定邮箱.
//  调用该函数前, 请确认:
//  1. 该用户存在并且 has_fixed
//  2. 该用户未当定邮箱
//  3. 该邮箱未绑定用户
func BindEmail(para *BindEmailParams) (err error) {
	parax := struct {
		*BindEmailParams
		AuthType int64 `sqlx:"auth_type"`
	}{
		BindEmailParams: para,
		AuthType:        AuthTypeEmail,
	}

	tx, err := db.GetDB().Beginx()
	if err != nil {
		return
	}

	// user_email 表增加一个 item
	stmt1, err := tx.Prepare("insert into user_email(user_id, nickname, email, has_bound) values(?, ?, ?, 1)")
	if err != nil {
		tx.Rollback()
		return
	}
	if _, err = stmt1.Exec(parax.UserId, parax.Email, parax.Email); err != nil {
		tx.Rollback()
		return
	}

	// user 更新 item
	stmt2, err := tx.PrepareNamed("update user set auth_types = auth_types|:auth_type where id=:user_id and has_fixed=1 and auth_types&:auth_type=0")
	if err != nil {
		tx.Rollback()
		return
	}
	rslt, err := stmt2.Exec(parax)
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
		err = fmt.Errorf("绑定邮箱 %s 到用户 %d 失败", para.Email, para.UserId)
		tx.Rollback()
		return
	}

	err = tx.Commit()
	return
}
