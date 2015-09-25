package model

import (
	"github.com/aiyi/go-user/db"
)

// 更新绑定邮箱
//  调用该函数前, 请确认:
//  1. 该用户存在, 并且 verified
//  2. 该用户已经绑定邮箱
func UpdateEmail(userId int64, email string) (err error) {
	para := struct {
		UserId   int64    `sqlx:"user_id"`
		Email    string   `sqlx:"email"`
		BindType BindType `sqlx:"bind_type"`
	}{
		UserId:   userId,
		Email:    email,
		BindType: BindTypeEmail,
	}

	stmt, err := db.GetNamedStmt("update user as A, user_email as B set B.email=:email where A.id=:user_id and A.verified=1 and A.bind_types&:bind_type<>0 and B.user_id=A.id and B.verified=1")
	if err != nil {
		return
	}
	_, err = stmt.Exec(para)
	return
}

// 更新绑定手机
//  调用该函数前, 请确认:
//  1. 该用户存在, 并且 verified
//  2. 该用户已经绑定手机
func UpdatePhone(userId int64, phone string) (err error) {
	para := struct {
		UserId   int64    `sqlx:"user_id"`
		Phone    string   `sqlx:"phone"`
		BindType BindType `sqlx:"bind_type"`
	}{
		UserId:   userId,
		Phone:    phone,
		BindType: BindTypePhone,
	}

	stmt, err := db.GetNamedStmt("update user as A, user_phone as B set B.phone=:phone where A.id=:user_id and A.verified=1 and A.bind_types&:bind_type<>0 and B.user_id=A.id and B.verified=1")
	if err != nil {
		return
	}
	_, err = stmt.Exec(para)
	return
}

// 更新绑定QQ
//  调用该函数前, 请确认:
//  1. 该用户存在, 并且 verified
//  2. 该用户已经绑定QQ
func UpdateQQ(userId int64, openid string) (err error) {
	para := struct {
		UserId   int64    `sqlx:"user_id"`
		OpenId   string   `sqlx:"openid"`
		BindType BindType `sqlx:"bind_type"`
	}{
		UserId:   userId,
		OpenId:   openid,
		BindType: BindTypeQQ,
	}

	stmt, err := db.GetNamedStmt("update user as A, user_qq as B set B.openid=:openid where A.id=:user_id and A.verified=1 and A.bind_types&:bind_type<>0 and B.user_id=A.id and B.verified=1")
	if err != nil {
		return
	}
	_, err = stmt.Exec(para)
	return
}

// 更新绑定微信
//  调用该函数前, 请确认:
//  1. 该用户存在, 并且 verified
//  2. 该用户已经绑定微信
func UpdateWechat(userId int64, openid string) (err error) {
	para := struct {
		UserId   int64    `sqlx:"user_id"`
		OpenId   string   `sqlx:"openid"`
		BindType BindType `sqlx:"bind_type"`
	}{
		UserId:   userId,
		OpenId:   openid,
		BindType: BindTypeWechat,
	}

	stmt, err := db.GetNamedStmt("update user as A, user_wechat as B set B.openid=:openid where A.id=:user_id and A.verified=1 and A.bind_types&:bind_type<>0 and B.user_id=A.id and B.verified=1")
	if err != nil {
		return
	}
	_, err = stmt.Exec(para)
	return
}

// 更新绑定微博
//  调用该函数前, 请确认:
//  1. 该用户存在, 并且 verified
//  2. 该用户已经绑定微博
func UpdateWeibo(userId int64, openid string) (err error) {
	para := struct {
		UserId   int64    `sqlx:"user_id"`
		OpenId   string   `sqlx:"openid"`
		BindType BindType `sqlx:"bind_type"`
	}{
		UserId:   userId,
		OpenId:   openid,
		BindType: BindTypeWeibo,
	}

	stmt, err := db.GetNamedStmt("update user as A, user_weibo as B set B.openid=:openid where A.id=:user_id and A.verified=1 and A.bind_types&:bind_type<>0 and B.user_id=A.id and B.verified=1")
	if err != nil {
		return
	}
	_, err = stmt.Exec(para)
	return
}
