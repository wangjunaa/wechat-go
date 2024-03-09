package dao

import (
	"github.com/spf13/viper"
	"log"
)

// 初始化yaml配置文件
func initConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("initConfig: 配置文件读取错误!")
		return
	}
	log.Println("initConfig Success")
}
