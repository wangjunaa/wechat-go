package redisLock

import (
	"demo/dao"
	"log"
	"time"
)

type RedisMux struct {
	Key               string           //锁名
	Id                string           //上锁的id，用于判断是否本机
	Expiration        time.Duration    //过期时间
	WatchDogCheckTime time.Duration    //看门狗检查时间，0则不开启
	Done              chan interface{} //看门狗结束指令
}

func (mux *RedisMux) Lock() {
	defer func() {
		if mux.WatchDogCheckTime != 0 {
			go mux.WatchDog()
		}
	}()
	lockSuccess := false
	var err error
	for !lockSuccess {
		lockSuccess, err = dao.Rdb.SetNX(dao.BgCtx, mux.Key, mux.Id, mux.Expiration).Result()
		if err != nil {
			log.Println("RedisLock.Lock:", err)
			return
		}
	}
}

func (mux *RedisMux) UnLock() {
	defer func() {
		if mux.WatchDogCheckTime != 0 {
			mux.Done <- "unlock"
		}
	}()
	id, err := dao.Rdb.Get(dao.BgCtx, mux.Key).Result()
	if err != nil {
		log.Println("RedisLock.Unlock:", err)
		return
	}
	if id != mux.Id {
		return
	}
	dao.Rdb.Del(dao.BgCtx, mux.Key)
}
func (mux *RedisMux) WatchDog() {
	for {
		time.Sleep(mux.WatchDogCheckTime)
		select {
		case <-mux.Done:
			return
		default:
			id, err := dao.Rdb.Get(dao.BgCtx, mux.Key).Result()
			if err != nil {
				return
			}
			if id == mux.Id {
				dao.Rdb.SetEX(dao.BgCtx, mux.Key, mux.Id, mux.Expiration)
			} else {
				log.Println("utils.redisLock.WatchDog: 锁被更改，但看门狗未结束")
				return
			}
		}
	}
}
