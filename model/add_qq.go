package model

import (
	"time"

	"github.com/aiyi/go-user/db"
	"github.com/aiyi/go-user/userid"
)

func AddQQ(openid, nickname string, timestamp int64) (userId int64, err error) {
	userId, err = userid.GetId()
	if err != nil {
		return
	}

	if timestamp == 0 {
		timestamp = time.Now().Unix()
	}

	para := struct {
		UserId     int64  `sqlx:"user_id"`
		AuthType   int64  `sqlx:"auth_type"`
		OpenId     string `sqlx:"openid"`
		Nickname   string `sqlx:"nickname"`
		Password   []byte `sqlx:"password"`
		Salt       []byte `sqlx:"salt"`
		CreateTime int64  `sqlx:"create_time"`
	}{
		UserId:     userId,
		AuthType:   AuthTypeQQ,
		OpenId:     openid,
		Nickname:   nickname,
		Password:   emptyByteSlice,
		Salt:       emptyByteSlice,
		CreateTime: timestamp,
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
	if _, err = stmt1.Exec(para.UserId, para.Nickname, para.OpenId); err != nil {
		tx.Rollback()
		return
	}

	// user 表增加一个 item
	stmt2, err := tx.PrepareNamed("insert into user(id, auth_types, password, salt, create_time, has_fixed) values(:user_id, :auth_type, :password, :salt, :create_time, 0)")
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
