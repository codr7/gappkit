package main

import (
	"gappkit/demo/bookng/core"
	"log"
)

func main() {
	log.Printf("Welcome to Bookng v%v", core.Version)
	db := core.NewDB("db")

	if err := db.Open(); err != nil {
		log.Fatal(err)
	}
	
	defer db.Close()
	
	r := db.NewResource()
	r.Name = "foo"

	if err := db.StoreResource(r); err != nil {
		log.Fatal(err)
	}

	lr, err := db.LoadResource(r.Id())

	if err != nil {
		log.Fatal(err)
	}
	
	log.Printf("%v\n", lr)
}
