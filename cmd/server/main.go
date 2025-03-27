package main

import (
	"chat/initialize"
	"chat/internal/global"
	r "chat/internal/handler/router"
)

func main() {
	initialize.SetUpViper()
	initialize.SetupLogger()
	initialize.SetupDatabase()

	router := r.Router()
	global.Logger.Info("服务正在启动......")
	if err := router.Run("localhost:8080"); err != nil {
		global.Logger.Fatal("服务启动失败" + err.Error())
	}
}
