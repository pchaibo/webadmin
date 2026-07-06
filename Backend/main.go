package main

import (
	"webadmin/binan"
	"webadmin/controller"
	"webadmin/router"
	"webadmin/task"
)

func main() {
	//go task.Sitestatus()
	go task.StartCli() //shell
	go bainstart()
	go controller.Broadcaster() //websocket

	router.Start()
}

func bainstart() {
	go binan.BinanPrice() //wss
	go binan.Taskuser()   //定时用户和币种价格
	go binan.Task()       //定时合约
}
