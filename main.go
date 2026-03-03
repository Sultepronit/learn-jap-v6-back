package main

import (
	"japv6/db"
	"japv6/server"
	"log"
)

func main() {
	log.Println("hello!")
	db.Open()
	db.Edit()
	server.Start()
}
