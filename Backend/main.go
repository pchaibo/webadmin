package main

import (
	"webadmin/router"
	"webadmin/task"
)

func main() {
	//go task.Sitestatus()
	go task.AddMin()
	go task.AddMax()
	router.Start()
}
