package utils

import (
	"chat/common/message"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

type Transfer struct {
	Conn net.Conn
	Buf [8094]byte
}

func (this *Transfer) Read()(mess message.Message,err error){
	// 先接收长度
	_,err = this.Conn.Read(this.Buf[:4])
	if err != nil {
		fmt.Println("接收长度失败")
		return
	}
	pkglen := binary.BigEndian.Uint32(this.Buf[:4])

	n,err := this.Conn.Read(this.Buf[:pkglen])

	if n != int(pkglen) {
		fmt.Println("长度不一致，发生丢包")
		err = errors.New("长度不一致，发生丢包")
	}
	if err != nil {
		return
	}

	err = json.Unmarshal(this.Buf[:pkglen],&mess)
	if err != nil {
		fmt.Println("err=",err)
	}
	return
}

func (this *Transfer) Send(data []byte)(err error) {

	var pkglen uint32
	pkglen = uint32(len(data))
	binary.BigEndian.PutUint32(this.Buf[:4],pkglen)

	//发送长度
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
	if  err != nil {
		fmt.Println("err=",err)
		return
	}

	return
}