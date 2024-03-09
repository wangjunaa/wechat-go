package dao

import (
	Model "demo/models"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func dbGetUserDsn() string {
	ip := viper.GetString("mysql.ip")
	user := viper.GetString("mysql.user")
	password := viper.GetString("mysql.password")
	port := viper.GetString("mysql.port")
	dbName := viper.GetString("mysql.dbName")
	charset := viper.GetString("mysql.charset")

	dsn := user + ":" + password + "@tcp(" + ip + ":" + port + ")/" + dbName + "?charset=" + charset
	return dsn
}
func initDB() {
	var err error
	dsn := dbGetUserDsn()
	log.Println("dsn:", dsn)
	DB, err = gorm.Open(mysql.Open(dsn))
	if err != nil {
		log.Println(err)
		return
	}
	err = DB.AutoMigrate(
		&Model.UserBasic{},
		&Model.GroupBasic{},
		&Model.Message{},
		&Model.FriendShip{},
	)
	if err != nil {
		log.Println("utils.dao.initDB:", err)
		return
	}

	log.Println("initDB Success")
}
