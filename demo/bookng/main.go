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
	root.Open()
	defer root.Close()
	
	var r Resource
	r.Id = resources.NextId()
	r.Name = "foo"
	resources.Store(r.Id, r)
}
