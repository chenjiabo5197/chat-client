package process

import (
	"common"
	"encoding/json"
	"fmt"
)

//展示服务器转发的群发聊天信息
func showGroupSms(mes *common.Message) (err error) {
	var smsMes common.SmsRespMes
	err = json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("smsResMes反序列化失败,err=", err)
		return
	}

	fmt.Printf("用户%s,向全体用户发送%s\n", smsMes.UserName, smsMes.Content)
	return
}

//展示服务器转发的1对1聊天信息
func showOne2OneSms(mes *common.Message, mesType common.MesType) (err error) {
	var smsMes common.SmsRespMes
	err = json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("smsResMes反序列化失败,err=", err)
		return
	}
	if mesType == common.RecvSmsToOneMesType {
		fmt.Printf("用户%s,向你单独发送%s\n", smsMes.SmsRespFrom, smsMes.Content)
	} else if mesType == common.OfflineRecvSmsToOneMesType {
		fmt.Printf("用户%s,离线期间向你单独发送%s\n", smsMes.SmsRespFrom, smsMes.Content)
	}
	return
}
