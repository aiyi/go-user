package model

import (
	"database/sql"
	"encoding/json"
	"strconv"

	"github.com/bradfitz/gomemcache/memcache"

	"github.com/aiyi/go-user/db"
	"github.com/aiyi/go-user/mc"
)

// user 基本信息
type User struct {
	Id          int64  `json:"id"           sqlx:"id"`
	Nickname    string `json:"nickname"     sqlx:"nickname"`
	BindTypes   int64  `json:"bind_types"   sqlx:"bind_types"`
	Password    []byte `json:"password"     sqlx:"password"`
	PasswordTag string `json:"password_tag" sqlx:"password_tag"`
	Salt        []byte `json:"salt"         sqlx:"salt"`
	CreateTime  int64  `json:"create_time"  sqlx:"create_time"`
	Verified    bool   `json:"verified"     sqlx:"verified"`
}

// memcache里user的key生成函数
func mcUserKey(userId int64) string {
	return "user/" + strconv.FormatInt(userId, 16)
}

// 从缓存里获取 user 信息, 如果没有找到返回 ErrNotFound.
func getUserFromCache(userId int64, user *User) (err error) {
	item, err := mc.Client().Get(mcUserKey(userId))
	if err != nil {
		if err == memcache.ErrCacheMiss {
			err = ErrNotFound
		}
		return
	}
	return json.Unmarshal(item.Value, user)
}

// 将 user 信息存入到缓存, 如果存在则更新.
func putUserToCache(user *User) (err error) {
	userBytes, err := json.Marshal(user)
	if err != nil {
		return
	}
	mcItem := memcache.Item{
		Key:   mcUserKey(user.Id),
		Value: userBytes,
		//Expiration: 300,
	}
	return mc.Client().Set(&mcItem)
}

func GetUser(userId int64) (user *User, err error) {
	user = &User{}

	if err = getUserFromCache(userId, user); err == nil {
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

	err = putUserToCache(user)
	return
}

func syncUserToCache(userId int64) (err error) {
	stmt, err := db.GetStmt("select * from user where id=?")
	if err != nil {
		return
	}

	var user User
	if err = stmt.Get(&user, userId); err != nil {
		if err == sql.ErrNoRows {
			err = ErrNotFound
		}
		return
	}

	return putUserToCache(&user)
}
