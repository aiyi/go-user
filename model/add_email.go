package model

import (
	"time"

	"github.com/aiyi/go-user/db"
	"github.com/aiyi/go-user/userid"
)

func AddEmail(email string, password, salt []byte, timestamp int64) (userId int64, err error) {
	userId, err = userid.GetId()
	if err != nil {
		return
	}

	if timestamp == 0 {
		timestamp = time.Now().Unix()
	}

	para := struct {
		UserId     int64    `sqlx:"user_id"`
		AuthType   AuthType `sqlx:"auth_type"`
		Email      string   `sqlx:"email"`
		Password   []byte   `sqlx:"password"`
		Salt       []byte   `sqlx:"salt"`
		CreateTime int64    `sqlx:"create_time"`
	}{
		UserId:     userId,
		AuthType:   AuthTypeEmail,
		Email:      email,
		Password:   password,
		Salt:       salt,
		CreateTime: timestamp,
	}

	tx, err := db.GetDB().Beginx()
	if err != nil {
		return
	}

	// user_email 表增加一个 item
	stmt1, err := tx.Prepare("insert into user_email(user_id, nickname, email, verified) values(?, ?, ?, 0)")
	if err != nil {
		tx.Rollback()
		return
	}
	if _, err = stmt1.Exec(para.UserId, para.Email, para.Email); err != nil {
		tx.Rollback()
		return
	}

	// user 表增加一个 item
	stmt2, err := tx.PrepareNamed("insert into user(id, auth_types, password, salt, create_time, verified) values(:user_id, :auth_type, :password, :salt, :create_time, 0)")
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
