package main

import (
	"example/task_manager/data"
	"example/task_manager/router"
)

func main() {
	data.InitMongoDB()
	r := router.SetupRouter()
	r.Run(":8080")
}