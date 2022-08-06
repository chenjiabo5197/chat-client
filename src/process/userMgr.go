package process

import (
	"common"
	"fmt"
	"model"
)

//客户端要维护一个map
var (
	onlineUsers map[int]*common.User = make(map[int]*common.User) //保存在线用户信息
	curUser     model.CurUser                                     //初始化，用于向服务器发送聊天信息
)

/*
	展示在线的用户列表
*/
func showOnlineUsers() {
	fmt.Println("在线用户列表：")
	for k := range onlineUsers {
		fmt.Println("用户id：", k)
	}
}

/*
	处理服务器返回的NotifyUserStatusMes
*/

func updateUserStatus(notifyMes *common.NotifyUserStatusMes) {

	user, ok := onlineUsers[notifyMes.UserId]
	if !ok {
		user = &common.User{
			UserId: notifyMes.UserId,
		}
	}
	if notifyMes.UserStatus == common.UserOnline {
		user.UserStatus = common.UserOnline
		onlineUsers[user.UserId] = user
	} else {
		delete(onlineUsers, user.UserId)
	}

	showOnlineUsers()
}
