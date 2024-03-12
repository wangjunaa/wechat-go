package dao

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"log"
)

func initRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.address"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
		PoolSize:     viper.GetInt("redis.poolSize"),
	})
	result, err := Rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
		return
	}
	log.Println("initRedis:", result)
}
