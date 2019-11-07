package controller

import (
	"chat/server/utils"
	"chat/common/message"
	"chat/server/model"
	"encoding/json"
	"fmt"
	"net"
)

type ServerUser struct {
	Conn net.Conn
	UserId int
}


func (this * ServerUser) Login (msg message.Message){
	var login message.LoginMes

	var mess message.Message
	var longinRes message.LoginResMes

	mess.Type = message.LoginResMesType

	err := json.Unmarshal([]byte(msg.Data),&login)
	if err != nil {
		fmt.Println("err=",err)
		return
	}
	err = model.MyUserModel.VerifyLogin(login.UserId,login.UserPwd)

	if err != nil {
		longinRes.Code = 401
		longinRes.Error = err.Error()
	}else{
		longinRes.Code = 200
		//在线用户
		this.UserId = login.UserId
		UserMger.Add(this)
		// 通知其他人
		this.NotifyOtherOfLogin(login.UserId)
		for id,_ := range UserMger.OnlineUsers {
			longinRes.UsersId = append(longinRes.UsersId,id)
		}

	}

	jsonStr,err := json.Marshal(longinRes)
	if err != nil {
		fmt.Println("err=",err)
	}

	mess.Data = string(jsonStr)

	data,err := json.Marshal(mess)
	tf := &utils.ServerTransfer{
		Conn:this.Conn,
	}

	err = tf.Send(data)
	if err != nil {
		fmt.Println("err=",err)
	}

	return

}

func (this * ServerUser) Register (msg message.Message){

	var register model.UserInfo
	err := json.Unmarshal([]byte(msg.Data),&register)

	if err != nil {
		fmt.Println("err=",err)
		return
	}

	var mess message.Message
	mess.Type = message.RegisterResMesType
	var respon message.RegisterResMes

	err = model.MyUserModel.Register(&register)
	if err != nil {
		respon.Code = 401
		respon.Error = err.Error()
	}else{
		respon.Code = 200
	}

	jsonStr,err := json.Marshal(respon)
	if err != nil {
		fmt.Println("err=",err)
	}

	mess.Data = string(jsonStr)

	data,err := json.Marshal(mess)
	tf := &utils.ServerTransfer{
		Conn:this.Conn,
	}

	err = tf.Send(data)
	if err != nil {
		fmt.Println("err=",err)
	}

	return

}

/**
	通知上线
 */
func (this *ServerUser) NotifyOtherOfLogin(userId int){
	onlineUsers := UserMger.All()
	for _,v := range onlineUsers{
		if v.UserId != userId {
			v.DoNotify(userId)
		}
	}
}
/**
	通知下线
 */
func (this *ServerUser) NotifyOtherOfOut(userId int){
	onlineUsers := UserMger.All()
	for _,v := range onlineUsers{
		if v.UserId != userId {
			v.DoOutNotify(userId)
		}
	}
}

func (this *ServerUser) DoNotify(userId int){
	tf := &utils.ServerTransfer{
		Conn: this.Conn,
	}

	var msg message.Message
	msg.Type = message.NotifyUserStatusMesType

	var notify message.NotifyUserStatusMes
	notify.UserId = userId
	notify.Status = 0

	jsonStr,err := json.Marshal(notify)
	if err != nil {
		fmt.Println("err=",err)
		return
	}

	msg.Data = string(jsonStr)
	data,err := json.Marshal(msg)
	if err != nil {
		fmt.Println("err=",err)
		return
	}

	err = tf.Send(data)
	if err != nil {
		fmt.Println("err=",err)
		return
	}
	return
}

func (this *ServerUser) DoOutNotify(userId int){
	tf := &utils.ServerTransfer{
		Conn: this.Conn,
	}

	var msg message.Message
	msg.Type = message.NotifyUserOutMesType

	var notify message.NotifyUserOutMes
	notify.UserId = userId
	notify.Status = 0

	jsonStr,err := json.Marshal(notify)
	if err != nil {
		fmt.Println("err=",err)
		return
	}

	msg.Data = string(jsonStr)
	data,err := json.Marshal(msg)
	if err != nil {
		fmt.Println("err=",err)
		return
	}

	err = tf.Send(data)
	if err != nil {
		fmt.Println("err=",err)
		return
	}
	return
}