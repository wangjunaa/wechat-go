package initSys

func Init() {
	initConfig()
	initDB()
	initRedis()
	initJwt()
}
