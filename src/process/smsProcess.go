package process

import (
	"common"
	"encoding/json"
	"fmt"
	"net"
	"utils"
)

// SmsProcessor 处理聊天消息的结构体
type SmsProcessor struct {
	Conn     net.Conn
	UserName string
	UserId   string
}

//向服务器发送广播聊天消息
func (sm *SmsProcessor) sendGroupSms(content string) (err error) {
	var smsMes common.SmsMes
	smsMes.Content = content
	smsMes.UserId = sm.UserId
	smsMes.UserName = sm.UserName
	data, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("smsMes序列化失败")
		return
	}

	err = sm.singleSendSms(data, common.SmsMesType)
	return
}

//向服务器发送1对1聊天消息
func (sm *SmsProcessor) sendToOneSms(content, userName string) (err error) {
	var smsMes common.SmsMes
	smsMes.Content = content
	smsMes.UserId = sm.UserId
	smsMes.UserName = sm.UserName
	smsMes.SmsMesTarget = userName
	data, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("smsMes序列化失败")
		return
	}

	err = sm.singleSendSms(data, common.SmsToOneMesType)
	return
}

// 传入要发送的data内容和type，然后组装发送到服务器
func (sm *SmsProcessor) singleSendSms(data []byte, mesType common.MesType) (err error) {
	var mes common.Message
	mes.Type = mesType
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
