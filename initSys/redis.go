package initSys

import (
	"demo/config"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"log"
)

func initRedis() {
	config.Rdb = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.address"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
		PoolSize:     viper.GetInt("redis.poolSize"),
	})
	result, err := config.Rdb.Ping(config.BgCtx).Result()
	if err != nil {
		log.Println("initRedis:", err)
		return
	}
	log.Println("initRedis:", result)
}
