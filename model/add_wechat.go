package model

import (
	"time"

	"github.com/chanxuehong/util/random"

	"github.com/aiyi/go-user/db"
	"github.com/aiyi/go-user/userid"
)

func AddByWechat(openid, nickname string, timestamp int64) (userId int64, err error) {
	userId, err = userid.GetId()
	if err != nil {
		return
	}

	if nickname == "" {
		nickname = openid
	}
	if timestamp == 0 {
		timestamp = time.Now().Unix()
	}

	para := struct {
		UserId      int64    `sqlx:"user_id"`
		BindType    BindType `sqlx:"bind_type"`
		OpenId      string   `sqlx:"openid"`
		Nickname    string   `sqlx:"nickname"`
		Password    []byte   `sqlx:"password"`
		PasswordTag []byte   `sqlx:"password_tag"`
		Salt        []byte   `sqlx:"salt"`
		CreateTime  int64    `sqlx:"create_time"`
	}{
		UserId:      userId,
		BindType:    BindTypeWechat,
		OpenId:      openid,
		Nickname:    nickname,
		Password:    emptyByteSlice,
		PasswordTag: random.NewRandomEx(),
		Salt:        emptyByteSlice,
		CreateTime:  timestamp,
	}

	tx, err := db.GetDB().Beginx()
	if err != nil {
		return
	}

	// user_wechat 表增加一个 item
	stmt1, err := tx.Prepare("insert into user_wechat(user_id, nickname, openid, verified) values(?, ?, ?, 0)")
	if err != nil {
		tx.Rollback()
		return
	}
	if _, err = stmt1.Exec(para.UserId, para.Nickname, para.OpenId); err != nil {
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
