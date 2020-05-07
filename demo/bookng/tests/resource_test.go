package tests

import (
	"gappkit/demo/bookng/core"
	"testing"
)

var db *core.DB

func fail(t *testing.T, err error) {
	t.Fatalf("%+v", err)
}

func setup(t *testing.T) {
	db = core.NewDB("testdb")

	if err := db.Drop(); err != nil {
		fail(t, err)
	}

	if err := db.Open(); err != nil {
		fail(t, err)
	}
}

func teardown(t *testing.T) {
	if err := db.Close(); err != nil {
		fail(t, err)
	}
}

func TestInitQuantity(t *testing.T) {
	setup(t)
	defer teardown(t)
	
	r := db.NewResource()
	r.Name = "foo"
	
	if err := r.Store(); err != nil {
		fail(t, err)
	}

	i := db.NewItem()
	i.Resource = r.Id()

	if err := i.Store(); err != nil {
		fail(t, err)
	}

	i = db.NewItem()
	i.Resource = r.Id()

	if err := i.Store(); err != nil {
		fail(t, err)
	}
}

func TestUpdateQuantity(t *testing.T) {
	setup(t)
	defer teardown(t)
	
	r := db.NewResource()
	r.Name = "foo"
	
	if err := r.Store(); err != nil {
		fail(t, err)
	}

	/*
	lr, err := db.LoadResource(r.Id())

	if err != nil {
		fail(t, err)
	}*/
}
