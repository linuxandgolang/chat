package controller

import (
	"chat/client/utils"
	"chat/common/message"
	"encoding/json"
	"fmt"
	"net"
)

type ServerSms struct {
	Conn net.Conn
}

func (this * ServerSms) SendGroupMessage(msg message.Message){
	var sms message.SmsMes

	err := json.Unmarshal([]byte(msg.Data),&sms)

	if err != nil {
		fmt.Println("err=",err)
	}

	data,err := json.Marshal(msg)

	if err != nil {
		fmt.Println("err=",err)
	}
	// 转发消息
	users := UserMger.All()

	for id,user := range users{
		if id == sms.UserId{
			continue
		}
		this.SendMesToEachOnlineUser(data,user.Conn)
	}

}

func (this *ServerSms) SendMesToEachOnlineUser(data []byte,conn net.Conn){
	tf := &utils.Transfer{
		Conn:conn,
	}

	err := tf.Send(data)

	if err != nil {
		fmt.Println("err=",err)
	}
	return
}
