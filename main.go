package main

import (
	"demo/dao"
	"demo/router"
	"log"
)

func main() {
	dao.Init()
	r := router.Router()
	err := r.Run(":8080")
	if err != nil {
		log.Println(err)
		return
	}

}
