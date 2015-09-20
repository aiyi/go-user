package model

import (
	"github.com/aiyi/go-user/db"
)

type AddWechatParams struct {
	UserId     int64  `sqlx:"user_id"`
	OpenId     string `sqlx:"openid"`
	Nickname   string `sqlx:"nickname"`
	CreateTime int64  `sqlx:"create_time"`
}

func AddWechat(para *AddWechatParams) (err error) {
	parax := struct {
		*AddWechatParams
		AuthType int64  `sqlx:"auth_type"`
		Password []byte `sqlx:"password"`
		Salt     []byte `sqlx:"salt"`
	}{
		AddWechatParams: para,
		AuthType:        AuthTypeWechat,
		Password:        emptyByteSlice,
		Salt:            emptyByteSlice,
	}

	tx, err := db.GetDB().Beginx()
	if err != nil {
		return
	}

	// user_wechat 表增加一个 item
	stmt1, err := tx.Prepare("insert into user_wechat(user_id, nickname, openid, has_bound) values(?, ?, ?, 0)")
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
