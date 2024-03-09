package dao

import (
	"github.com/spf13/viper"
	"log"
)

func initJwt() {
	SecretKey = viper.GetString("token.secretKey")
	Issuer = viper.GetString("token.Issuer")
	ExpiresTime = viper.GetInt("token.ExpiresTime")
	log.Println("intJwt:", SecretKey)
}
