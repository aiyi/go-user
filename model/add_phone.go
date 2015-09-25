package model

import (
	"time"

	"github.com/chanxuehong/util/random"

	"github.com/aiyi/go-user/db"
	"github.com/aiyi/go-user/userid"
)

// password, salt 可以为 nil
func AddPhone(phone string, password, salt []byte, timestamp int64) (userId int64, err error) {
	userId, err = userid.GetId()
	if err != nil {
		return
	}

	if password == nil {
		password = emptyByteSlice
	}
	if salt == nil {
		salt = emptyByteSlice
	}
	if timestamp == 0 {
		timestamp = time.Now().Unix()
	}

	para := struct {
		UserId      int64    `sqlx:"user_id"`
		BindType    BindType `sqlx:"bind_type"`
		Phone       string   `sqlx:"phone"`
		Password    []byte   `sqlx:"password"`
		PasswordTag []byte   `sqlx:"password_tag"`
		Salt        []byte   `sqlx:"salt"`
		CreateTime  int64    `sqlx:"create_time"`
	}{
		UserId:      userId,
		BindType:    BindTypePhone,
		Phone:       phone,
		Password:    password,
		PasswordTag: random.NewRandomEx(),
		Salt:        salt,
		CreateTime:  timestamp,
	}

	tx, err := db.GetDB().Beginx()
	if err != nil {
		return
	}

	// user_phone 表增加一个 item
	stmt1, err := tx.Prepare("insert into user_phone(user_id, nickname, phone, verified) values(?, ?, ?, 0)")
	if err != nil {
		tx.Rollback()
		return
	}
	if _, err = stmt1.Exec(para.UserId, para.Phone, para.Phone); err != nil {
		tx.Rollback()
		return
	}

	// user 表增加一个 item
	stmt2, err := tx.PrepareNamed("insert into user(id, bind_types, password, password_tag, salt, create_time, verified) values(:user_id, :bind_type, :password, :password_tag, :salt, :create_time, 0)")
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
