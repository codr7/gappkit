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

	d := dom.NewNode("html")
	h := d.NewNode("head")
	h.NewNode("title").Append("Title")
	b := d.NewNode("body")
	b.NewNode("a").Set("href", "https://foo.bar").Append("Foobar")
	d.Write(os.Stdout)
	
	defer db.Close()
}
