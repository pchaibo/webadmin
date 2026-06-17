package main

import (
	router "webadmin/router"
	"webadmin/task"
)

func main() {
	go task.Sitestatus()
	router.Start()
}
