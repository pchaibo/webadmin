package main

import (
	"webadmin/binan"
	"webadmin/router"
	"webadmin/task"
)

func main() {
	//go task.Sitestatus()
	go task.StartCli()
	go binan.BinanPrice()
	router.Start()
}
