package handle

import (
	"demo/tools/redisLock"
	"time"
)

var (
	//redis hash存储key
	hUserKey = "user:1"

	//redis分布式锁
	rUserMux = redisLock.RedisMux{
		Key:               "redisUsers",
		Id:                "1",
		Expiration:        10 * time.Second,
		WatchDogCheckTime: 0,
		Done:              make(chan interface{}),
	}
	//sql分布式锁
	sUserMux = redisLock.RedisMux{
		Key:               "sqlUsers",
		Id:                "1",
		Expiration:        10 * time.Second,
		WatchDogCheckTime: 0,
		Done:              make(chan interface{}),
	}
	sGroupMux = redisLock.RedisMux{
		Key:               "sqlGroup",
		Id:                "1",
		Expiration:        10 * time.Second,
		WatchDogCheckTime: 0,
		Done:              make(chan interface{}),
	}
	sFriendMux = redisLock.RedisMux{
		Key:               "sqlFriend",
		Id:                "1",
		Expiration:        10 * time.Second,
		WatchDogCheckTime: 0,
		Done:              make(chan interface{}),
	}
)
