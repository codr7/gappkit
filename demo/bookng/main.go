package main

import (
	"gappkit/demo/bookng/core"
	"log"
)

func main() {
	log.Printf("Welcome to Bookng v%v", core.Version)
	db := core.NewDB("db")

	if err := db.Open(); err != nil {
		log.Fatalf("%+v", err)
	}
	
	defer db.Close()
}
