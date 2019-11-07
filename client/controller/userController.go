package controller

import (
	"chat/client/utils"
	"chat/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type UserController struct {

}
type UserInfo struct {
	UserId int `json:"user_id"`
	UserName string `json:"user_name"`
	UserPwd string `json:"user_pwd"`
	UserStatus int `json:"user_status"`
}

type CurrUser struct {
	Conn net.Conn
	UserId int
}

func (this * UserController) Login(userId int,userPwd string)  {

	var msg message.Message
	msg.Type = message.LoginMesType

	var loginInfo message.LoginMes

	loginInfo.UserId = userId
	loginInfo.UserPwd = userPwd

	jsonStr,err := json.Marshal(loginInfo)
	if err != nil {
		fmt.Println("err=",err)
	}

	msg.Data = string(jsonStr)

	data,err := json.Marshal(msg)
	if err != nil {
		fmt.Println("err=",err)
	}
	// 发送消息
	conn,err := net.Dial("tcp",":8888")
	if err != nil {
		fmt.Println("err=",err)
	}
	defer conn.Close()
	// 发送长度
	var pkglen uint32
	pkglen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:4],pkglen)
	n,err := conn.Write(buf[:4])

	if n != 4 || err != nil {
		fmt.Println("err=",err)
	}
	//fmt.Printf("客户端，发送消息的长度=%d 内容=%s \n", len(data), string(data))
	// 发送消息本身
	_,err = conn.Write(data)
	if err != nil {
		fmt.Println("err=",err)
	}

	// 接收消息
	tf := &utils.Transfer{
		Conn:conn,
	}
	msg,err = tf.Read()
	if err != nil {
		fmt.Println("err=",err)
	}
	var info message.LoginResMes
	err = json.Unmarshal([]byte(msg.Data),&info)

	if err != nil {
		fmt.Println("err=",err)
	}

	if info.Code == 200 {

		CurUser.Conn = conn
		CurUser.UserId = userId

		fmt.Println("登录成功，当前在线用户列表如下:")
		for _, v := range info.UsersId {
			//如果我们要求不显示自己在线,下面我们增加一个代码
			if v == userId {
				fmt.Println("▶自己->id:", v)
			}else{
				fmt.Println("用户id:", v)
			}
			user := &message.User{
				UserId: v,
				Status: 0,
			}
			ClientOnlineUsers[v]= user
		}
		fmt.Print("\n")
		//这里我们还需要在客户端启动一个协程
		//该协程保持和服务器端的通讯.如果服务器有数据推送给客户端
		//则接收并显示在客户端的终端.
		go serverProcessMes(conn)

		for{
			showMenu()
		}
	}else {
		fmt.Println("err=",info.Error)
	}

}

func (this * UserController) Register (userId int,userPwd string,userName string)  {
	var mess message.Message
	mess.Type = message.RegisterMesType

	var register message.RegisterMes

	register.UserId = userId
	register.UserName = userName
	register.UserPwd = userPwd

	jsonStr,err := json.Marshal(register)

	if err != nil {
		fmt.Println("err=",err)
		return
	}

	mess.Data = string(jsonStr)

	data,err := json.Marshal(mess)

	if err != nil {
		fmt.Println("err=",err)
		return
	}
	// 发送消息
	conn,err := net.Dial("tcp",":8888")
	if err != nil {
		fmt.Println("err=",err)
	}
	defer conn.Close()
	tf := &utils.Transfer{
		Conn: conn,
	}
	fmt.Println(string(data))
	err = tf.Send(data)
	if err != nil {
		fmt.Println("err=",err)
		return
	}

	returns,err := tf.Read()
	if err != nil {
		fmt.Println("err=",err)
		return
	}

	var info message.RegisterResMes

	err = json.Unmarshal([]byte(returns.Data),&info)

	if err != nil {
		fmt.Println("err=",err)
	}

	if info.Code == 200 {
		fmt.Println("注册成功，请重新登录")
		os.Exit(0)
	}else{
		fmt.Println(info.Error)
		os.Exit(0)
	}
	return



}
