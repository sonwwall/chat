package global

import (
	config "chat/configs"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Config     *config.Config
	Logger     *zap.Logger
	Mysql      *gorm.DB
	Redis      *redis.Client
	GinCOntext *gin.Context
)
