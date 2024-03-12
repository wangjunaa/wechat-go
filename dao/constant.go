package dao

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// SecretKey 令牌密钥
var SecretKey string
var Issuer string
var ExpiresTime int

// DB mysql数据库
var DB *gorm.DB

// Rdb redis数据库
var Rdb *redis.Client
