package main

import (
	"demo/initSys"
	"demo/router"
	"log"
)

func main() {
	initSys.Init()
	r := router.Router()
	err := r.Run(":8080")
	if err != nil {
		log.Println(err)
		return
	}
}
