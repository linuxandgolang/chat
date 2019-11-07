package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type UserInfo struct {
	UserId int `json:"user_id"`
	UserName string `json:"user_name"`
	UserPwd string `json:"user_pwd"`
	UserStatus int `json:"user_status"`
}

var MyUserModel *UserModel

type UserModel struct {
	pool *redis.Pool
}

func NewUserModel(pool *redis.Pool)(user *UserModel){
	user = &UserModel{
		pool:pool,
	}
	return
}

func (this *UserModel) VerifyLogin(userId int,userPwd string)(err error){

	conn := this.pool.Get()
	defer conn.Close()

	res,err := redis.String(conn.Do("hget","users",userId))
	if err != nil {
		fmt.Println("err=",err)
		return
	}
	var user UserInfo
	err = json.Unmarshal([]byte(res),&user)
	if err != nil {
		fmt.Println("err=",err)
		return
	}
	if user.UserPwd != userPwd {
		err = errors.New("账号密码不一致")
		return

	}

	return
}
func (this * UserModel) Register(userInfo *UserInfo)(err error){

	conn := this.pool.Get()
	defer conn.Close()
	data,err := json.Marshal(userInfo)
	if err != nil {
		fmt.Println("err=",err)
	}
	ishas,err := redis.String(conn.Do("hget","users",userInfo.UserId))
	if ishas != "" {
		err = errors.New("账号已存在")
		return
	}
	res,err := redis.Bool(conn.Do("hset","users",userInfo.UserId,string(data)))
	if err != nil {
		fmt.Println("err=",err)
		return
	}
	if !res {
		err = errors.New("注册失败")
		return
	}

	return
}


