package controller

import (
	"chat/client/utils"
	"chat/common/message"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

func showMenu(){
	fmt.Println("-------1. 显示在线用户列表---------")
	fmt.Println("-------2. 发送消息---------")
	fmt.Println("-------3. 退出系统---------")
	fmt.Println("请选择(1-3):")
	var key int
	var content string

	fmt.Scanf("%d\n", &key)
	switch key {
	case 1:
		outputOnlineUsers()
	case 2:
		fmt.Println("你想对大家说的什么:)")
		fmt.Scanf("%s\n", &content)
		sendMsg := &SendMsg{
			
		};
		err := sendMsg.SendToOthers(content)
		if err != nil {
			fmt.Println("发送失败，err=",err)
		}
	case 3:
		fmt.Println("你选择退出了系统...")
		os.Exit(0)
	default :
		fmt.Println("你输入的选项不正确..")
	}
}
func serverProcessMes(conn net.Conn)  {
	tf := &utils.Transfer{
		Conn:conn,
	}

	for{
		msg,err := tf.Read()
		if err != nil {
			fmt.Println("err=",err)
		}

		switch msg.Type {
			case message.NotifyUserStatusMesType:
				// 有人上线了
				var notifyUserStatusMes message.NotifyUserStatusMes
					err = json.Unmarshal([]byte(msg.Data), &notifyUserStatusMes)
				//2. 把这个用户的信息，状态保存到客户map[int]User中
				updateUserStatus(&notifyUserStatusMes)
			case message.NotifyUserOutMesType:
				var notifyUserOutMes message.NotifyUserOutMes
				err = json.Unmarshal([]byte(msg.Data), &notifyUserOutMes)
				updateUserOut(&notifyUserOutMes)
			case message.SmsMesType:
				// 群发消息
				var smsMes message.SmsMes
				err = json.Unmarshal([]byte(msg.Data),&smsMes)
				outputGroupMsg(smsMes)
		}
	}
}
