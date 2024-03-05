package initSys

import (
	"demo/config"
	"github.com/spf13/viper"
	"log"
)

func initJwt() {
	config.SecretKey = viper.GetString("token.secretKey")
	config.Issuer = viper.GetString("token.Issuer")
	config.ExpiresTime = viper.GetInt("token.ExpiresTime")
	log.Println("intJwt:", config.SecretKey)
}
