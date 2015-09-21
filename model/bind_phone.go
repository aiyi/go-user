package model

import (
	"fmt"

	"github.com/aiyi/go-user/db"
)

type BindPhoneParams struct {
	UserId int64  `sqlx:"user_id"` // 绑定到这个用户
	Phone  string `sqlx:"phone"`
}

// 给用户绑定手机.
//  调用该函数前, 请确认:
//  1. 该用户存在并且 has_fixed
//  2. 该用户未绑定手机
//  3. 该手机未绑定用户
func BindPhone(para *BindPhoneParams) (err error) {
	parax := struct {
		*BindPhoneParams
		AuthType int64 `sqlx:"auth_type"`
	}{
		BindPhoneParams: para,
		AuthType:        AuthTypePhone,
	}

	tx, err := db.GetDB().Beginx()
	if err != nil {
		return
	}

	// user_phone 表增加一个 item
	stmt1, err := tx.Prepare("insert into user_phone(user_id, nickname, phone, has_fixed) values(?, ?, ?, 1)")
	if err != nil {
		tx.Rollback()
		return
	}
	if _, err = stmt1.Exec(parax.UserId, parax.Phone, parax.Phone); err != nil {
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
		err = fmt.Errorf("绑定手机 %s 到用户 %d 失败", para.Phone, para.UserId)
		tx.Rollback()
		return
	}

	err = tx.Commit()
	return
}
