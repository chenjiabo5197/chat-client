package model

import (
	"common"
	"net"
)

/*
	保存当前登录用户的信息，包含conn连接和user的信息，用于向服务端发送消息
*/
type CurUser struct {
	Conn net.Conn
	common.User
}
