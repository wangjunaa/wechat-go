package handler

import (
	"demo/utils/redisLock"
	"time"
)

var (
	//redis hash存储key
	hUserKey = "user:1"

	//redis分布式锁
	userMux = redisLock.RedisMux{
		Key:               "redisUsers",
		Id:                "1",
		Expiration:        10 * time.Second,
		WatchDogCheckTime: 0,
		Done:              make(chan interface{}),
	}
	groupMux = redisLock.RedisMux{
		Key:               "sqlGroup",
		Id:                "1",
		Expiration:        10 * time.Second,
		WatchDogCheckTime: 0,
		Done:              make(chan interface{}),
	}

	friendMux = redisLock.RedisMux{
		Key:               "sqlFriend",
		Id:                "1",
		Expiration:        10 * time.Second,
		WatchDogCheckTime: 0,
		Done:              make(chan interface{}),
	}
)
