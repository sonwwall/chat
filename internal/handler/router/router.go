package router

import (
	"chat/internal/controllers"
	"chat/internal/global"
	"chat/internal/handler/http"
	"chat/internal/middleware"
	"chat/internal/repository"
	"chat/internal/service"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()

	userRepo := repository.NewUserRepository(global.Mysql)
	userService := service.NewUserService(userRepo)
	userHandler := http.NewUserHandler(userService)

	userGroup := r.Group("/user")
	{
		userGroup.POST("/register", userHandler.UserRegister)
		userGroup.POST("/login", userHandler.UserLogin)
	}

	messageGroup := r.Group("/chat")
	messageGroup.Use(middleware.JwtAuthMiddleware())
	ms := service.NewMessageService(global.Redis)
	chatController := controllers.NewChatController(global.Mysql, ms)
	messageGroup.GET("/rooms/:room_id/connect", chatController.ConnectWebSocket)
	messageGroup.POST("/messages", chatController.SendMessage)

	return r
}
