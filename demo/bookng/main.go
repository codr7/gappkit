package main

import (
	"gappkit/demo/bookng/core"
	"gappkit/db"
	"log"
)

type Resource struct {
	Id db.RecordId
	Name string
}

func main() {
	log.Printf("Welcome to Bookng v%v", bookng.Version)

	root := db.NewRoot("db")
	resources := root.NewTable("resource")
	resources.NewColumn("Name")

	if err := root.Open(); err != nil {
		log.Fatal(err)
	}
	
	defer root.Close()
	
	var r Resource
	r.Id = resources.NextId()
	r.Name = "foo"

	if err := resources.Store(r.Id, &r); err != nil {
		log.Fatal(err)
	}

	var lr Resource

	if err := resources.Load(r.Id, &lr); err != nil {
		log.Fatal(err)
	}
	
	log.Printf("%v\n", lr)
}
