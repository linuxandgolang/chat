package controller

import "fmt"

var UserMger *UserManager

type UserManager struct {
	OnlineUsers map[int]*ServerUser
}



func (this *UserManager) Add(user *ServerUser){
	UserMger.OnlineUsers[user.UserId] = user
}

func (this *UserManager) Del (userId int){
	delete(UserMger.OnlineUsers,userId)
}

func (this *UserManager) All()map[int]*ServerUser{
	return UserMger.OnlineUsers
}

func (this *UserManager) GetOnlineUserById(userId int)(user *ServerUser,err error){

	user,ok := UserMger.OnlineUsers[userId]
	if !ok {
		err = fmt.Errorf("用户%d 不存在", userId)
		return
	}
	return
}
