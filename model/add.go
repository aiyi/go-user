package model

import (
	"github.com/aiyi/go-user/db"
	timex "github.com/chanxuehong/util/time"
)

const (
	AuthTypeEmail         = 1 << iota // 邮箱-密码
	AuthTypePhonePassword             // 手机-密码
	AuthTypePhoneCode                 // 手机-短信验证码
	AuthTypeWechat                    // 微信
	AuthTypeQQ                        // QQ
	AuthTypeWeibo                     // 微博
)

var emptyByteSlice = []byte{}

type AddByEmailParams struct {
	AuthType   int64      `json:"auth_type"`
	Email      string     `json:"email"`
	Password   []byte     `json:"password"`
	Salt       []byte     `json:"salt"`
	CreateTime timex.Time `json:"create_time"`
}

func AddByEmail(para *AddByEmailParams) (userid int64, err error) {
	para.AuthType = AuthTypeEmail

	tx, err := db.GetDB().Beginx()
	if err != nil {
		return
	}

	// user 表增加一个 item
	stmt, err := tx.PrepareNamed("insert into user(auth_types, password, salt, create_time) values(:auth_type, :password, :salt, :create_time)")
	if err != nil {
		tx.Rollback()
		return
	}

	rslt, err := stmt.Exec(para)
	if err != nil {
		tx.Rollback()
		return
	}

	userid, err = rslt.LastInsertId()
	if err != nil {
		tx.Rollback()
		return
	}

	// user_email 表增加一个 item
	stmt2, err := tx.Prepare("insert into user_email(userid, email) values(?, ?)")
	if err != nil {
		tx.Rollback()
		return
	}

	if _, err = stmt2.Exec(userid, para.Email); err != nil {
		tx.Rollback()
		return
	}

	err = tx.Commit()
	return
}
