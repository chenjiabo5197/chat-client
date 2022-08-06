package process

import (
	"common"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"model"
	"net"
	"utils"
)

type UserProcessor struct {
}

/*
	完成登录函数，传入userId和userPwd,返回服务器返回的登陆结果
*/
func (up *UserProcessor) Login(userId int, userPwd string) (err error) {
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
	loginMes.UserId = userId
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

	//为了确保tcp发送消息的准确性，先发送mes的长度给服务器，再发送mes消息本体
	//先获取data这个byte切片的长度，然后将长度数据转化为一个byte切片
	var dataLen uint32
	dataLen = uint32(len(data))
	var dataLenbytes [4]byte
	//此函数可以将传入的一个uint32的数值转化为一个byte切片
	binary.BigEndian.PutUint32(dataLenbytes[:4], dataLen)
	//发送长度数据
	n, err := conn.Write(dataLenbytes[:4])
	if n != 4 || err != nil {
		fmt.Println("发送data长度失败,err=", err)
		return
	}

	// fmt.Printf("发送数据长度=%d, 内容=%s\n", len(data), string(data))

	//发送数据
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("发送data失败,err=", err)
		return
	}

	/*
		用程序休眠来检测刚开始服务器端一直打印错误的原因，程序在此暂停，发现是因为客户端关闭
		服务器端的conn.Read函数不阻塞，导致一直读取数据，但是读取不到，所以报错
		后面在服务器端 增加err == io.EOF判断是否是客户端关闭了conn连接
	*/
	// time.Sleep(10 * time.Second)
	tf := utils.Transfer{
		Conn: conn,
	}
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("读取服务器返回出错,err=", err)
		return
	}

	//将服务器返回消息反序列化为LoginResMes结构体
	var loginResMes common.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if err != nil {
		fmt.Println("反序列化服务器返回结果出错,err=", err)
		return
	}

	//判断是否登录成功
	if loginResMes.ResCode == 200 {
		//初始化curUser
		curUser = model.CurUser{}
		curUser.Conn = conn
		curUser.UserId = userId
		curUser.UserStatus = common.UserOnline

		//显示服务器返回的在线用户列表
		for _, v := range loginResMes.UsersId {
			if v == userId {
				continue
			}
			//初始化onlineUsers
			user := &common.User{
				UserId:     v,
				UserStatus: common.UserOnline,
			}
			onlineUsers[v] = user

			fmt.Println("在线用户 : ", v)
		}
		go serverProcessMes(conn)
		showMenu()
		return nil
	} else {
		return errors.New(loginResMes.Error)
	}
}

/*
	完成注册函数，传入userId和userPwd和userName,返回服务器返回的注册结果
*/
func (up *UserProcessor) Register(user common.User) (err error) {

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

	//为了确保tcp发送消息的准确性，先发送mes的长度给服务器，再发送mes消息本体
	//先获取data这个byte切片的长度，然后将长度数据转化为一个byte切片
	var dataLen uint32
	dataLen = uint32(len(data))
	var dataLenbytes [4]byte
	//此函数可以将传入的一个uint32的数值转化为一个byte切片
	binary.BigEndian.PutUint32(dataLenbytes[:4], dataLen)
	//发送长度数据
	n, err := conn.Write(dataLenbytes[:4])
	if n != 4 || err != nil {
		fmt.Println("发送data长度失败,err=", err)
		return
	}

	//发送数据
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("发送data失败,err=", err)
		return
	}

	tf := utils.Transfer{
		Conn: conn,
	}
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("读取服务器返回出错,err=", err)
		return
	}

	//将服务器返回消息反序列化为LoginResMes结构体
	var rigisterResMes common.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &rigisterResMes)
	if err != nil {
		fmt.Println("反序列化服务器返回结果出错,err=", err)
		return
	}

	//判断是否注册成功
	if rigisterResMes.ResCode == 200 {
		return nil
	} else {
		return errors.New(rigisterResMes.Error)
	}
}
