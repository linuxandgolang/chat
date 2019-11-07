package utils

import (
	"chat/common/message"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
)

type ServerTransfer struct {
	Conn net.Conn
	Buf [8096]byte
}
func (this * ServerTransfer) Send(data []byte)(err error){

	// 先发送长度
	var pkglen uint32
	pkglen = uint32(len(data))
	binary.BigEndian.PutUint32(this.Buf[0:4],pkglen)

	n,err := this.Conn.Write(this.Buf[:4])

	if n != 4 {
		err = errors.New("发送长度时丢包")
		return
	}
	if err != nil {
		fmt.Println("err=",err)
		return
	}

	// 发送本身

	n,err = this.Conn.Write(data)

	if n != int(pkglen) {
		err = errors.New("发送长度时丢包")
		return
	}
	if err != nil {
		fmt.Println("err=",err)
		return
	}
	return

}


func (this *ServerTransfer) Read()(msg message.Message,err error){

	_,err  = this.Conn.Read(this.Buf[:4])
	if err != nil {
		if err == io.EOF {
			fmt.Println("客户端退出，服务器端也退出..")
			return
		} else {
			fmt.Println("接收长度失败err=",err)
		}
		return
	}

	pkglen := binary.BigEndian.Uint32(this.Buf[:4])

	n,err := this.Conn.Read(this.Buf[:pkglen])
	if err != nil {
		fmt.Println("err=",err)
		return
	}
	if n != int(pkglen) {
		fmt.Println("长度不一致，发生丢包")
		err = errors.New("长度不一致，发生丢包")
		return
	}
	err = json.Unmarshal(this.Buf[:pkglen], &msg)
	if err != nil {
		fmt.Println("json.Unmarsha err=", err)
		return
	}
	return
}
