package initialize

import (
	"chat/internal/global"
	"github.com/spf13/viper"
)

func SetUpViper() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.SetConfigFile("./configs/config.yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic("Read config failed " + err.Error())
	}

	err = viper.Unmarshal(&global.Config)
	if err != nil {
		panic("Unmarshal config failed " + err.Error())
	}

}
