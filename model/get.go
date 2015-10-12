package model

import (
	"database/sql"

	"github.com/aiyi/go-user/db"
)

// user 基本信息
type User struct {
	Id          int64    `json:"id"           sqlx:"id"`
	Nickname    string   `json:"nickname"     sqlx:"nickname"`
	BindTypes   BindType `json:"bind_types"   sqlx:"bind_types"`
	Password    []byte   `json:"-"            sqlx:"password"`
	PasswordTag string   `json:"-"            sqlx:"password_tag"`
	Salt        []byte   `json:"-"            sqlx:"salt"`
	CreateTime  int64    `json:"create_time"  sqlx:"create_time"`
	Verified    bool     `json:"verified"     sqlx:"verified"`
}

func Get(userId int64) (user *User, err error) {
	user = &User{}

	if err = getFromCache(userId, user); err == nil {
		return
	}
	if err != ErrNotFound {
		return
	}

	// 缓存没有找到, 从数据库读取
	stmt, err := db.GetStmt("select * from user where id=?")
	if err != nil {
		return
	}
	if err = stmt.Get(user, userId); err != nil {
		if err == sql.ErrNoRows {
			err = ErrNotFound
		}
		return
	}

	err = putToCache(user)
	return
}

func GetByEmail(email string) (user *User, err error) {
	stmt, err := db.GetStmt("select A.id, A.nickname, A.bind_types, A.password, A.password_tag, A.salt, A.create_time, A.verified from user as A, user_email as B where B.email=? and A.id=B.user_id and A.verified=B.verified and A.bind_types&?<>0")
	if err != nil {
		return
	}

	user = &User{}
	if err = stmt.Get(user, email, BindTypeEmail); err != nil {
		if err == sql.ErrNoRows {
			err = ErrNotFound
		}
		return
	}
	return
}

func GetByPhone(phone string) (user *User, err error) {
	stmt, err := db.GetStmt("select A.id, A.nickname, A.bind_types, A.password, A.password_tag, A.salt, A.create_time, A.verified from user as A, user_phone as B where B.phone=? and A.id=B.user_id and A.verified=B.verified and A.bind_types&?<>0")
	if err != nil {
		return
	}

	user = &User{}
	if err = stmt.Get(user, phone, BindTypePhone); err != nil {
		if err == sql.ErrNoRows {
			err = ErrNotFound
		}
		return
	}
	return
}

func GetByQQ(openid string) (user *User, err error) {
	stmt, err := db.GetStmt("select A.id, A.nickname, A.bind_types, A.password, A.password_tag, A.salt, A.create_time, A.verified from user as A, user_qq as B where B.openid=? and A.id=B.user_id and A.verified=B.verified and A.bind_types&?<>0")
	if err != nil {
		return
	}

	user = &User{}
	if err = stmt.Get(user, openid, BindTypeQQ); err != nil {
		if err == sql.ErrNoRows {
			err = ErrNotFound
		}
		return
	}
	return
}

func GetByWechat(openid string) (user *User, err error) {
	stmt, err := db.GetStmt("select A.id, A.nickname, A.bind_types, A.password, A.password_tag, A.salt, A.create_time, A.verified from user as A, user_wechat as B where B.openid=? and A.id=B.user_id and A.verified=B.verified and A.bind_types&?<>0")
	if err != nil {
		return
	}

	user = &User{}
	if err = stmt.Get(user, openid, BindTypeWechat); err != nil {
		if err == sql.ErrNoRows {
			err = ErrNotFound
		}
		return
	}
	return
}

func GetByWeibo(openid string) (user *User, err error) {
	stmt, err := db.GetStmt("select A.id, A.nickname, A.bind_types, A.password, A.password_tag, A.salt, A.create_time, A.verified from user as A, user_weibo as B where B.openid=? and A.id=B.user_id and A.verified=B.verified and A.bind_types&?<>0")
	if err != nil {
		return
	}

	user = &User{}
	if err = stmt.Get(user, openid, BindTypeWeibo); err != nil {
		if err == sql.ErrNoRows {
			err = ErrNotFound
		}
		return
	}
	return
}
