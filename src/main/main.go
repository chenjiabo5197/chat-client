package main

import (
	"common"
	"fmt"
	"process"
	"strconv"
)

func Init() (recommend int) {

	fmt.Println("-------------------欢迎登录多人聊天系统-----------------")
	fmt.Println("                     1、登录聊天系统                   ")
	fmt.Println("                     2、注册用户                       ")
	fmt.Println("                     3、退出系统                       ")
	fmt.Println("请选择（1-3）：")
	fmt.Println("------------------------------------------------------")

	var selectInput string
	_, _ = fmt.Scanln(&selectInput)

	recommend, err := strconv.Atoi(selectInput)
	if err != nil || recommend <= 0 || recommend > 3 {
		fmt.Println("输入错误，请重新输入")
		return 0
	} else {
		return recommend
	}
}

/*
	用于接收登录用户输入的用户id和密码
*/
func inputUser() (userId int, userPwd string, err error) {

	fmt.Println("请输入用户id:")
	var inputId string
	_, _ = fmt.Scanln(&inputId)

	userId, err = strconv.Atoi(inputId)
	if err != nil {
		fmt.Println("用户id必须是数字")
		return
	}
	fmt.Println("请输入用户密码:")
	_, _ = fmt.Scanln(&(userPwd))
	return
}

/*
	用于接收注册用户输入的用户id，密码和名字
*/
func inputRegisterUser() (user common.User, err error) {

	fmt.Println("请输入用户id:")
	var inputId string
	_, _ = fmt.Scanln(&inputId)
	user = common.User{}

	user.UserId, err = strconv.Atoi(inputId)
	if err != nil {
		fmt.Println("用户id必须是数字")
		return
	}
	fmt.Println("请输入用户密码:")
	_, _ = fmt.Scanln(&(user.UserPwd))
	fmt.Println("请输入用户姓名:")
	_, _ = fmt.Scanln(&(user.UserName))
	return
}

func main() {

	//用于控制循环的变量
	var loop = true
	//记录用户输入的变量
	var recommend int

	for loop {
		recommend = Init()
		if recommend == 0 {
			continue
		} else {
			switch recommend {
			case 1:
				userId, userPwd, err := inputUser()
				if err != nil {
					continue
				}
				up := process.UserProcessor{}
				err = up.LoginHandler(userId, userPwd)
				if err != nil {
					fmt.Println("登陆失败,err=", err)
				}
				loop = false
			case 2:
				user, err := inputRegisterUser()
				if err != nil {
					continue
				}
				up := process.UserProcessor{}
				err = up.RegisterHandler(user)
				if err != nil {
					fmt.Println("注册失败,err=", err)
				} else {
					fmt.Println("注册成功,请重新登陆")
				}
				loop = false
			case 3:
				fmt.Println("退出系统")
				loop = false
			default:
				fmt.Println("error")
			}
		}
	}
}
