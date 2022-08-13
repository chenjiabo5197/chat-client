package process

import (
	"common"
	"encoding/json"
	"fmt"
	"model"
	"utils"
)

//客户端要维护一个map
var (
	onlineUsers = make(map[int]*common.User) //保存在线用户信息,有缺陷，要下线
	//curUser     model.CurUser                //初始化，用于向服务器发送聊天信息
)

//展示在线的用户列表
func queryAllOnlineUsers(user *model.CurUser) error {
	//fmt.Println("在线用户列表：")
	//for _, user := range onlineUsers {
	//	fmt.Printf("%s在线\n", user.UserName)
	//}
	// 定义查询消息结构体,将自身信息传入，服务器过滤自身信息
	mes := common.Message{}
	mes.Type = common.QueryAllOnlineType
	User := &common.User{
		UserId: user.UserId,
	}
	userByte, _ := json.Marshal(User)
	mes.Data = string(userByte)
	mesByte, _ := json.Marshal(mes)
	tf := utils.Transfer{
		Conn: user.Conn,
	}
	err := tf.WritePkg(mesByte)
	if err != nil {
		fmt.Printf("send query online mes err, err=%s\n", err.Error())
		return err
	}
	return nil
}

func showAllOnlineUser(onlineData *string) {
	var onlineUsername []string
	err := json.Unmarshal([]byte(*onlineData), &onlineUsername)
	if err != nil {
		fmt.Printf("unmarshal onlineUsername err, err=%s\n", err.Error())
		return
	}
	fmt.Println("在线用户：")
	for i := range onlineUsername {
		fmt.Println(onlineUsername[i])
	}
}

/*
	处理服务器返回的NotifyUserStatusMes
*/

//func updateUserStatus(notifyMes *common.NotifyUserStatusMes) {
//
//	user, ok := onlineUsers[notifyMes.UserId]
//	if !ok {
//		user = &common.User{
//			UserId:   notifyMes.UserId,
//			UserName: notifyMes.UserName,
//		}
//	}
//	if notifyMes.UserStatus == common.UserOnline {
//		user.UserStatus = common.UserOnline
//		onlineUsers[user.UserId] = user
//	} else {
//		delete(onlineUsers, user.UserId)
//	}
//
//	showOnlineUsers()
//}
