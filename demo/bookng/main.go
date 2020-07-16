package main

import (
	"gappkit/demo/bookng/pkg"
	"gappkit/dom"
	"log"
	"os"
	"time"
)

func main() {
	log.Printf("Welcome to Bookng v%v", bookng.Version)
	db := bookng.NewDB("db")

	if err := db.Open(time.Now()); err != nil {
		log.Fatalf("%+v", err)
	}

	defer db.Close()
}
