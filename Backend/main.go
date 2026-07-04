package main

import (
	"webadmin/binan"
	"webadmin/controller"
	"webadmin/router"
	"webadmin/task"
)

func main() {
	//go task.Sitestatus()
	go task.StartCli()
	go binan.BinanPrice()
	go controller.Broadcaster()
	go binan.Taskuser()
	router.Start()
}
