package process

import (
	"common"
	"encoding/json"
	"fmt"
	"net"
	"utils"
)

//处理聊天消息的结构体
type SmsProcessor struct {
	Conn     net.Conn
	UserName string
	UserId   int
}

/*
	向服务器发送聊天消息的方法
*/
func (sm *SmsProcessor) sendGroupSms(content string) (err error) {
	var smsMes common.SmsMes
	smsMes.Content = content
	smsMes.UserId = sm.UserId
	data, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("smsMes序列化失败")
		return
	}

	var mes common.Message
	mes.Type = common.SmsMesType
	mes.Data = string(data)
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("mes序列化失败")
		return
	}

	tf := utils.Transfer{
		Conn: sm.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("发送sms消息失败,err=", err)
		return
	}
	return

}
