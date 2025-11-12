package main

import (
	"task_manager/data"
	"task_manager/database"
	"task_manager/router"
)

func main() {
	database.Init()
	data.InitUserService()
	data.InitTaskService()

	r := router.SetupRouter()
	r.Run(":8080")
}
