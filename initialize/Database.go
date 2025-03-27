package initialize

import (
	"chat/internal/global"
	"chat/migrations"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetupDatabase() {
	SetupMysql()
	SetupRedis()
}

func SetupMysql() {
	mysqlConfig := global.Config.
		DatabaseConfig.MysqlConfig
	dsn := mysqlConfig.Username + ":" + mysqlConfig.
		Password + "@tcp(" + mysqlConfig.Addr + ")/" + mysqlConfig.
		DB + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		global.Logger.Fatal("数据库连接失败" + err.Error())
	}
	global.Mysql = db
	global.Logger.Info("数据库链接成功")

	migrations.Migrate(global.Mysql)
}

func SetupRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     global.Config.DatabaseConfig.RedisConfig.Addr,
		Password: global.Config.DatabaseConfig.RedisConfig.Password,
		DB:       global.Config.DatabaseConfig.RedisConfig.DB,
	})
	global.Redis = rdb
	global.Logger.Info("redis连接成功")
}
