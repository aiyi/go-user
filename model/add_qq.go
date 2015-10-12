package model

import (
	"time"

	"github.com/aiyi/go-user/db"
	"github.com/aiyi/go-user/userid"
)

// 通过 QQ 注册一个账户.
//  如果 nickname 为空, 则默认为 openid
//  如果 timestamp == 0 则默认使用当前时间
func AddByQQ(openid, nickname string, timestamp int64) (user *User, err error) {
	userId, err := userid.GetId()
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
		PasswordTag string   `sqlx:"password_tag"`
		Salt        []byte   `sqlx:"salt"`
		CreateTime  int64    `sqlx:"create_time"`
		Verified    bool     `sqlx:"verified"`
	}{
		UserId:      userId,
		BindType:    BindTypeQQ,
		OpenId:      openid,
		Nickname:    nickname,
		Password:    emptyByteSlice,
		PasswordTag: NewPasswordTag(),
		Salt:        emptyByteSlice,
		CreateTime:  timestamp,
		Verified:    defaultVerified,
	}

	tx, err := db.GetDB().Beginx()
	if err != nil {
		return
	}

	// user_qq 表增加一个 item
	stmt1, err := tx.Prepare("insert into user_qq(user_id, openid, verified) values(?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return
	}
	if _, err = stmt1.Exec(para.UserId, para.OpenId, para.Verified); err != nil {
		tx.Rollback()
		return
	}

	// user 表增加一个 item
	stmt2, err := tx.PrepareNamed("insert into user(id, nickname, bind_types, password, password_tag, salt, create_time, verified) values(:user_id, :nickname, :bind_type, :password, :password_tag, :salt, :create_time, :verified)")
	if err != nil {
		tx.Rollback()
		return
	}
	if _, err = stmt2.Exec(para); err != nil {
		tx.Rollback()
		return
	}

	if err = tx.Commit(); err != nil {
		return
	}

	user = &User{
		Id:          para.UserId,
		Nickname:    para.Nickname,
		BindTypes:   para.BindType,
		Password:    para.Password,
		PasswordTag: para.PasswordTag,
		Salt:        para.Salt,
		CreateTime:  para.CreateTime,
		Verified:    para.Verified,
	}
	return
}
