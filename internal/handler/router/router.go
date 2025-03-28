package router

import (
	"chat/internal/global"
	"chat/internal/handler/http"
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

	return r
}
