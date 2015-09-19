package model

import (
	"github.com/aiyi/go-user/db"
	"github.com/chanxuehong/util/timex"
)

type AddByEmailParams struct {
	UserId     int64      `sqlx:"id"`
	AuthType   int64      `sqlx:"auth_type"`
	Email      string     `sqlx:"email"`
	Password   []byte     `sqlx:"password"`
	Salt       []byte     `sqlx:"salt"`
	CreateTime timex.Time `sqlx:"create_time"`
}

func AddByEmail(para *AddByEmailParams) (err error) {
	para.AuthType = AuthTypeEmail

	tx, err := db.GetDB().Beginx()
	if err != nil {
		return
	}

	// user_email 表增加一个 item
	stmt1, err := tx.Prepare("insert into user_email(userid, email) values(?, ?)")
	if err != nil {
		tx.Rollback()
		return
	}

	if _, err = stmt1.Exec(para.UserId, para.Email); err != nil {
		tx.Rollback()
		return
	}

	// user 表增加一个 item
	stmt2, err := tx.PrepareNamed("insert into user(id, auth_types, password, salt, create_time) values(:id, :auth_type, :password, :salt, :create_time)")
	if err != nil {
		tx.Rollback()
		return
	}

	if _, err = stmt2.Exec(para); err != nil {
		tx.Rollback()
		return
	}

	err = tx.Commit()
	return
}
