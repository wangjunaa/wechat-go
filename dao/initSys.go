package dao

func Init() {
	initConfig()
	initDB()
	initRedis()
	initJwt()
}
