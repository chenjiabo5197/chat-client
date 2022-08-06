package process

import (
	"common"
	"encoding/json"
	"fmt"
)

/*
	展示服务器转发的聊天信息
*/
func showGroupSms(mes *common.Message) (err error) {
	var smsMes common.SmsResMes
	err = json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("smsResMes反序列化失败,err=", err)
		return
	}

	fmt.Printf("用户%d,发送%s\n", smsMes.UserId, smsMes.Content)
	return
}
