package model

import (
	"time"

	"github.com/aiyi/go-user/db"
	"github.com/aiyi/go-user/userid"
)

type AddQQParams struct {
	OpenId   string `sqlx:"openid"`
	Nickname string `sqlx:"nickname"`
}

func AddQQ(para *AddQQParams) (userId int64, err error) {
	userId, err = userid.GetId()
	if err != nil {
		return
	}

	parax := struct {
		*AddQQParams
		UserId     int64  `sqlx:"user_id"`
		AuthType   int64  `sqlx:"auth_type"`
		Password   []byte `sqlx:"password"`
		Salt       []byte `sqlx:"salt"`
		CreateTime int64  `sqlx:"create_time"`
	}{
		AddQQParams: para,
		UserId:      userId,
		AuthType:    AuthTypeQQ,
		Password:    emptyByteSlice,
		Salt:        emptyByteSlice,
		CreateTime:  time.Now().Unix(),
	}

	tx, err := db.GetDB().Beginx()
	if err != nil {
		return
	}

	// user_qq 表增加一个 item
	stmt1, err := tx.Prepare("insert into user_qq(user_id, nickname, openid, has_fixed) values(?, ?, ?, 0)")
	if err != nil {
		tx.Rollback()
		return
	}
	if _, err = stmt1.Exec(parax.UserId, parax.Nickname, parax.OpenId); err != nil {
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
