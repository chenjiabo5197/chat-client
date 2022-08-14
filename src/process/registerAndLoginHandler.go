package process

import (
	"common"
	"encoding/json"
	"errors"
	"fmt"
	"model"
	"net"
	"utils"
)

type UserProcessor struct {
}

// LoginHandler 完成登录函数，传入userName和userPwd,返回服务器返回的登陆结果
func (up *UserProcessor) LoginHandler(userName string, userPwd string) (err error) {
	// fmt.Printf("用户输入的id=%d,密码=%s", userId, userPwd)
	// return nil

	//连接到服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("连接服务器出错，err=", err)
		return
	}

	//在函数最后关闭conn
	defer conn.Close()

	//通过conn向服务器发送消息
	//先创建一个发送的消息实例,并指明是发送的登陆消息
	var mes common.Message
	mes.Type = common.LoginMesType
	//再创建一个登陆的消息实例
	var loginMes common.LoginMes
	loginMes.UserId = utils.GetMd5Value(userName)
	loginMes.UserName = userName
	loginMes.UserPwd = userPwd

	//将登陆的消息序列化,返回的是一个byte切片
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("登陆消息json Marshal出错")
		return
	}
	mes.Data = string(data)

	//将发送的消息序列化,若转换不出错，data即要发送的数据，data为切片类型
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("发送消息json Marshal出错")
		return
	}

	tf := utils.Transfer{
		Conn: conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Printf("send data to server err, err=%v\n", err)
		return err
	}

	/*
		用程序休眠来检测刚开始服务器端一直打印错误的原因，程序在此暂停，发现是因为客户端关闭
		服务器端的conn.Read函数不阻塞，导致一直读取数据，但是读取不到，所以报错
		后面在服务器端 增加err == io.EOF判断是否是客户端关闭了conn连接
	*/
	// time.Sleep(10 * time.Second)
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("读取服务器返回出错,err=", err)
		return
	}

	fmt.Printf("service resp login mes=%s\n", utils.Struct2String(mes))

	//将服务器返回消息反序列化为LoginResMes结构体
	var loginRespMes common.LoginRespMes
	err = json.Unmarshal([]byte(mes.Data), &loginRespMes)
	if err != nil {
		fmt.Println("反序列化服务器返回结果出错,err=", err)
		return
	}

	//判断是否登录成功
	if loginRespMes.RespCode == 200 {
		//初始化curUser
		curUser := &model.CurUser{}
		curUser.Conn = conn
		curUser.UserId = utils.GetMd5Value(userName)
		curUser.UserName = loginRespMes.UserName
		curUser.UserStatus = common.UserOnline

		//显示服务器返回的在线用户列表
		//for _, v := range loginRespMes.UsersId {
		//	if v == userId {
		//		continue
		//	}
		//	//初始化onlineUsers
		//	user := &common.User{
		//		UserId:     v,
		//		UserStatus: common.UserOnline,
		//	}
		//	onlineUsers[v] = user
		//
		//	fmt.Println("在线用户 : ", v)
		//}
		go serverProcessMes(conn)
		showMenu(curUser)
		return nil
	} else {
		return errors.New(loginRespMes.Error)
	}
}

/*
	完成注册函数，传入userId和userPwd和userName,返回服务器返回的注册结果
*/
func (up *UserProcessor) RegisterHandler(user common.User) (err error) {

	//连接到服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("连接服务器出错，err=", err)
		return
	}

	//在函数最后关闭conn
	defer conn.Close()

	//通过conn向服务器发送消息
	//先创建一个发送的消息实例,并指明是发送的注册消息
	var mes common.Message
	mes.Type = common.RegisterMesType

	rigister := common.RegisterMes{
		User: user,
	}

	//将注册的消息序列化,返回的是一个byte切片
	data, err := json.Marshal(rigister)
	if err != nil {
		fmt.Println("注册消息json Marshal出错")
		return
	}
	mes.Data = string(data)

	//将发送的消息序列化,若转换不出错，data即要发送的数据，data为切片类型
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("发送消息json Marshal出错")
		return
	}

	tf := utils.Transfer{
		Conn: conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Printf("send data to server err, err=%v\n", err)
		return err
	}

	tf = utils.Transfer{
		Conn: conn,
	}
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("读取服务器返回出错,err=", err)
		return
	}
	// 解析服务器的返回数据
	err = ParseServerResp(&mes)
	return
}
