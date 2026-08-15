package main

import (
	"webadmin/binan"
	"webadmin/controller"
	"webadmin/router"
	"webadmin/task"
)

func getshell() {
	//go task.Sitestatus()
	go task.StartCli()
}

func main_test() {
	//binan.TestSendSMS()
	binan.PostEmail()
}
func main() {

	//go getshell()
	go bainstart()
	go controller.Broadcaster() //websocket

	router.Start()
}

func bainstart() {
	go binan.BinanPrice() //wss
	go binan.Taskuser()   //定时用户和币种价格
	go binan.Task()       //定时合约
}
