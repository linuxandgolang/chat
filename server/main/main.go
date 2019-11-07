package main

import (
	"chat/common/message"
	"chat/server/controller"
	"chat/server/model"
	"chat/server/utils"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"io"
	"net"
	"time"
)
var pool *redis.Pool

func init(){
	initPool()
	initUserModel()
	initUserMgar()
}
func initPool(){
	pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp","localhost:6379")
		},
		MaxIdle:         16,
		MaxActive:       0,
		IdleTimeout:     300 * time.Second,
	}
}

func initUserModel(){
	model.MyUserModel = model.NewUserModel(pool)
}

func initUserMgar(){
	controller.UserMger = &controller.UserManager{
		OnlineUsers: make(map[int]*controller.ServerUser,1024),
	}
}

func main() {

	listen,err := net.Listen("tcp",":8888")
	if err != nil {
		fmt.Println("err=",err)
	}

	for{
		conn,err := listen.Accept()
		if err != nil {
			fmt.Println("err=",err)
		}

		fmt.Printf("%s:连接成功\n",conn.RemoteAddr())

		go router(conn)
	}
}

func router(conn net.Conn){
	defer conn.Close()

	for{
		tf := &utils.ServerTransfer{
			Conn: conn,
		}
		res,err := tf.Read()
		if err != nil {
			if err == io.EOF {
				for _,v :=range controller.UserMger.OnlineUsers{
					if v.Conn == conn{
						delete(controller.UserMger.OnlineUsers,v.UserId)
						//通知客户端 有链接下线
						server := &controller.ServerUser{
							Conn:   conn,
							UserId: v.UserId,
						}
						server.NotifyOtherOfOut(v.UserId)
					}
				}
			} else {
				fmt.Println("read err=", err)

			}
			return
		}
		switch res.Type {
			case message.LoginMesType:
				user := &controller.ServerUser{
					Conn:conn,
				}
				user.Login(res)
			case message.RegisterMesType:
				user := &controller.ServerUser{
					Conn:conn,
				}
				user.Register(res)
			case message.SmsMesType:
				sms := &controller.ServerSms{
					Conn:conn,
				}
				sms.SendGroupMessage(res)
			default:
				fmt.Println("类型错误")
		}
	}

}
