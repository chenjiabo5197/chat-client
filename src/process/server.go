package process

import (
	"common"
	"encoding/json"
	"fmt"
	"model"
	"net"
	"strconv"
	"utils"
)

/*
	1、显示登陆成功页面
	2、保持和服务器通信
	3、当读取到服务器发送的消息后，显示在界面
*/

func showMenu(user *model.CurUser) {

	recommend := true
	for recommend {
		fmt.Printf("------------恭喜%s登录系统----------\n", user.UserName)
		fmt.Println("          1、显示在线用户列表          ")
		fmt.Println("          2、发送消息          ")
		fmt.Println("          3、消息列表          ")
		fmt.Println("          4、退出系统          ")
		fmt.Println("请选择(1-4):")
		var input string
		var content string
		fmt.Scanln(&input)
		key, err := strconv.Atoi(input)
		if err != nil || key <= 0 || key > 4 {
			fmt.Println("请输入正确的选项")
		}
		switch key {
		case 1:
			// fmt.Println("1")
			queryAllOnlineUsers(user)
		case 2:
			//fmt.Println("2")
			fmt.Println("请输入要发送的消息:")
			fmt.Scanln(&content)
			sp := SmsProcessor{
				Conn:     user.Conn,
				UserName: user.UserName,
				UserId:   user.UserId,
			}
			err := sp.sendGroupSms(content)
			if err != nil {
				fmt.Println("发送消息失败,err=", err)
			}
		case 3:
			fmt.Println("3")
		case 4:
			fmt.Println("退出系统")
			recommend = false
		default:
			fmt.Println("error")

		}
	}

}

/*
	和服务器保持连接
*/
func serverProcessMes(conn net.Conn) {
	//创建一个Transfer对象，不停地读取服务器发送的消息
	tf := utils.Transfer{
		Conn: conn,
	}

	for {
		fmt.Println("客户端正在等待读取服务器消息")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("读取消息失败,err=", err)
			return
		}
		switch mes.Type {
		case common.NotifyUserStatusMesType:
			//上线/下线通知
			var notifyMes common.NotifyUserStatusMes
			err := json.Unmarshal([]byte(mes.Data), &notifyMes)
			if err != nil {
				fmt.Println("notifyMes反序列化失败,err=", err)
				return
			}
			fmt.Printf("recv server online and offline mes, data=%s\n", utils.Struct2String(notifyMes))
			//更新onlineUsers 这个map
			//updateUserStatus(&notifyMes)
		case common.SmsRespMesType:
			//服务器转发的消息
			err = showGroupSms(&mes)
			if err != nil {
				fmt.Println("err=", err)
				return
			}
		case common.AllOnlineRespType: // 用户在线消息
			showAllOnlineUser(&mes.Data)
		default:
			fmt.Println("未知数据类型")
		}

		//打印读取的消息
		//fmt.Println("mes=",mes)
	}
}
