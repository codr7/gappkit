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
	
	var r core.Resource
	r.Id = db.Resources.NextId()
	r.Name = "foo"
	
	if err := db.Resources.Store(r.Id, &r); err != nil {
		log.Fatal(err)
	}

	var lr core.Resource

	if err := db.Resources.Load(r.Id, &lr); err != nil {
		log.Fatal(err)
	}
	
	log.Printf("%v\n", lr)
}
