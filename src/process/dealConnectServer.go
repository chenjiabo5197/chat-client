package process

import (
	"common"
	"encoding/json"
	"errors"
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
		fmt.Println("          2、发送广播消息          ")
		fmt.Println("          3、发送1对1消息          ")
		fmt.Println("          4、消息列表          ")
		fmt.Println("          5、退出系统          ")
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
			ol := Online{}
			err := ol.queryAllOnlineUsers(user)
			if err != nil {
				return
			}
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
			var userName string
			fmt.Println("请输入要发送消息的对象名称:")
			fmt.Scanln(&userName)
			fmt.Println("请输入要发送消息的内容:")
			fmt.Scanln(&content)
			sp := SmsProcessor{
				Conn:     user.Conn,
				UserName: user.UserName,
				UserId:   user.UserId,
			}
			err := sp.sendToOneSms(content, userName)
			if err != nil {
				fmt.Println("发送消息失败,err=", err)
			}
		case 4:
			fmt.Println("3")
		case 5:
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
			case common.NotifyUserStatusMesType:  //上线/下线通知
				var notifyMes common.NotifyUserStatusMes
				err := json.Unmarshal([]byte(mes.Data), &notifyMes)
				if err != nil {
					fmt.Println("notifyMes反序列化失败,err=", err)
					return
				}
				fmt.Printf("recv server online and offline mes, data=%s\n", utils.Struct2String(notifyMes))
				//更新onlineUsers 这个map
				//updateUserStatus(&notifyMes)
			case common.RecvSmsMesType:  //服务器转发的群发消息
				err = showGroupSms(&mes)
				if err != nil {
					fmt.Println("err=", err)
					return
				}
			case common.RecvSmsToOneMesType:  //服务器转发的1对1消息
				err = showOne2OneSms(&mes)
				if err != nil {
					fmt.Println("err=", err)
					return
				}
			case common.SmsRespMesType:  //发送的1对1或群发消息失败还是成功
				err = ParseServerResp(&mes)
				if err != nil {
					fmt.Println("发送消息失败，err=", err)
				}else {
					fmt.Println("发送消息成功")
				}
			case common.AllOnlineRespType: // 用户在线消息
				ol := Online{}
				ol.showAllOnlineUser(&mes.Data)
			default:
				fmt.Println("未知数据类型")
			}
	}
}

func ParseServerResp(message *common.Message) (err error) {
	//将服务器返回消息反序列化为LoginResMes结构体
	var respMes common.StatusRespMes
	err = json.Unmarshal([]byte(message.Data), &respMes)
	if err != nil {
		return errors.New("反序列化服务器返回结果出错")
	}

	//判断返回消失成功还是失败
	if respMes.RespCode == 200 {
		return nil
	} else {
		return errors.New(respMes.Error)
	}
}
