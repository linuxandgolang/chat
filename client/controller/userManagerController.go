package controller

import (
	"chat/common/message"
	"fmt"
)

var ClientOnlineUsers map[int]*message.User = make(map[int]*message.User,1024)
var CurUser CurrUser

func updateUserStatus (notifyUserStatusMes *message.NotifyUserStatusMes){

	user := &message.User{
		UserId:notifyUserStatusMes.UserId,
		Status:notifyUserStatusMes.Status,
	}
	ClientOnlineUsers[notifyUserStatusMes.UserId] = user
	outputOnlineUsers()
}

func outputOnlineUsers(){
	fmt.Println("当前在线用户列表:")
	for id, _ := range ClientOnlineUsers{
		if id == CurUser.UserId {
			fmt.Println("▶自己id:", id)
		}else{
			fmt.Println("用户id:", id)
		}
	}
}

func updateUserOut(notifyUserOutMes *message.NotifyUserOutMes){
	delete(ClientOnlineUsers,notifyUserOutMes.UserId)
}


