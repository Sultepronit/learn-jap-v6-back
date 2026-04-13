package main

import (
	"japv6/db"
	"japv6/server"
	"log"
)

func main() {
	db.Open()
	// db.Edit()
	// v, err := db.GetVersion("word_cards")
	// log.Println("v:", v, "err:", err)
	log.Println("Testing!")
	server.Start()
}
