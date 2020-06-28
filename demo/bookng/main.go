package main

import (
	"gappkit/demo/bookng/core"
	"gappkit/dom"
	"log"
	"os"
)

func main() {
	log.Printf("Welcome to Bookng v%v", core.Version)
	db := core.NewDB("db")

	if err := db.Open(); err != nil {
		log.Fatalf("%+v", err)
	}

	d := dom.NewNode("html")
	h := d.NewNode("head")
	h.NewNode("title").Append("Title")
	b := d.NewNode("body")
	b.NewNode("a").Set("href", "https://foo.bar").Append("Foobar")
	d.Write(os.Stdout)
	
	defer db.Close()
}
