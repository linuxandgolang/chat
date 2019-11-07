package controller

import (
	"chat/client/utils"
	"chat/common/message"
	"encoding/json"
	"fmt"
)

type SendMsg struct {
	
}

func (this *SendMsg) SendToOthers(content string)(err error){
	
	var mess message.Message
	mess.Type = message.SmsMesType
	
	var sms message.SmsMes
	sms.UserId = CurUser.UserId
	sms.Content = content
	
	jsonStr,err := json.Marshal(sms)
	if err != nil {
		fmt.Println("err=",err)
		return
	}
	
	mess.Data = string(jsonStr)
	
	tf:= &utils.Transfer{
		Conn: CurUser.Conn,
	}
	data,err := json.Marshal(mess)
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

func outputGroupMsg(msg message.SmsMes){
	// 显示群发消息
	info := fmt.Sprintf("用户id:\t%d 对大家说:\t%s",msg.UserId,msg.Content)
	fmt.Println(info)
	fmt.Println()
}


