package process

import (
	"encoding/json"
	"fmt"
	"vs_code/project0007/client/utils"
	"vs_code/project0007/common"
)

//处理聊天消息的结构体
type SmsProcessor struct {

}

/*
	向服务器发送聊天消息的方法
*/
func (sm *SmsProcessor) sendGroupSms (content string) (err error) {
	var smsMes common.SmsMes
	smsMes.Content = content
	smsMes.UserId = curUser.UserId
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

	tf := utils.Transfer {
		Conn: curUser.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("发送sms消息失败,err=",err)
		return
	}
	return

}
