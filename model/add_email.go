package model

import (
	"time"

	"github.com/aiyi/go-user/db"
	"github.com/aiyi/go-user/userid"
)

type AddEmailParams struct {
	Email    string `sqlx:"email"`
	Password []byte `sqlx:"password"`
	Salt     []byte `sqlx:"salt"`
}

func AddEmail(para *AddEmailParams) (err error) {
	userId, err := userid.GetId()
	if err != nil {
		return
	}

	parax := struct {
		*AddEmailParams
		UserId     int64 `sqlx:"user_id"`
		AuthType   int64 `sqlx:"auth_type"`
		CreateTime int64 `sqlx:"create_time"`
	}{
		AddEmailParams: para,
		UserId:         userId,
		AuthType:       AuthTypeEmail,
		CreateTime:     time.Now().Unix(),
	}

	tx, err := db.GetDB().Beginx()
	if err != nil {
		return
	}

	// user_email 表增加一个 item
	stmt1, err := tx.Prepare("insert into user_email(user_id, nickname, email, has_bound) values(?, ?, ?, 0)")
	if err != nil {
		tx.Rollback()
		return
	}
	if _, err = stmt1.Exec(parax.UserId, parax.Email, parax.Email); err != nil {
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
