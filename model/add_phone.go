package model

import (
	"time"

	"github.com/aiyi/go-user/db"
	"github.com/aiyi/go-user/userid"
)

type AddPhoneParams struct {
	Phone    string `sqlx:"phone"`
	Password []byte `sqlx:"password"` // 可以为 nil
	Salt     []byte `sqlx:"salt"`     // 可以为 nil
}

func AddPhone(para *AddPhoneParams) (err error) {
	if para.Password == nil {
		para.Password = emptyByteSlice
	}
	if para.Salt == nil {
		para.Salt = emptyByteSlice
	}

	userId, err := userid.GetId()
	if err != nil {
		return
	}

	parax := struct {
		*AddPhoneParams
		UserId     int64 `sqlx:"user_id"`
		AuthType   int64 `sqlx:"auth_type"`
		CreateTime int64 `sqlx:"create_time"`
	}{
		AddPhoneParams: para,
		UserId:         userId,
		AuthType:       AuthTypePhone,
		CreateTime:     time.Now().Unix(),
	}

	tx, err := db.GetDB().Beginx()
	if err != nil {
		return
	}

	// user_phone 表增加一个 item
	stmt1, err := tx.Prepare("insert into user_phone(user_id, nickname, phone, has_bound) values(?, ?, ?, 0)")
	if err != nil {
		tx.Rollback()
		return
	}
	if _, err = stmt1.Exec(parax.UserId, parax.Phone, parax.Phone); err != nil {
		tx.Rollback()
		return
	}

	// user 表增加一个 item
	stmt2, err := tx.PrepareNamed("insert into user(id, auth_types, password, salt, create_time, has_fixed) values(:user_id, :auth_type, :password, :salt, :create_time, 0)")
	if err != nil {
		tx.Rollback()
		return
	}
	if _, err = stmt2.Exec(parax); err != nil {
		tx.Rollback()
		return
	}

	err = tx.Commit()
	return
}
