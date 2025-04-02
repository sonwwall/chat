package main

import (
	"chat/initialize"
	"chat/internal/global"
	r "chat/internal/handler/router"
	"chat/internal/handler/ws"
)

func main() {
	initialize.SetUpViper()
	initialize.SetupLogger()
	initialize.SetupDatabase()

	router := r.Router()

	go ws.HubInstance.Run()

	global.Logger.Info("服务正在启动......")
	if err := router.Run("0.0.0.0:8080"); err != nil {
		global.Logger.Fatal("服务启动失败" + err.Error())
	}
}
