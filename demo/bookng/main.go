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
	
	r := db.NewResource()
	r.Name = "foo"

	if err := r.Store(); err != nil {
		log.Fatalf("%+v", err)
	}

	lr, err := db.LoadResource(r.Id())

	if err != nil {
		log.Fatal("%+v", err)
	}
	
	log.Printf("%v\n", lr)
}
