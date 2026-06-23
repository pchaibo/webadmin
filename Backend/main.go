package main

import (
	"webadmin/router"
	"webadmin/task"
)

func main() {
	//go task.Sitestatus()
	go task.StartCli()
	router.Start()
}
