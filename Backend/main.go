package main

import (
	"webadmin/binan"
	"webadmin/ccxt"
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
	ccxt.Tdemo()

}
func main_bak() {
	ccxt.Tdemo()
	ccxt.FetchAccount()
}
func main() {
	go bainstart()
	go controller.Broadcaster() //websocket

	router.Start()
}

func bainstart() {
	go binan.BinanPrice() //wss
	go binan.Taskuser()   //定时用户和币种价格
	go binan.Task()       //定时合约
}
