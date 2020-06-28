package tests

import (
	"gappkit/demo/bookng/core"
	"github.com/pkg/errors"
	"testing"
	"time"
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

	if err := db.Open(time.Now()); err != nil {
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

	if q, err := r.AvailableQuantity(r.StartTime, i.StartTime); err != nil {
		fail(t, err)
	} else if q != 1 {
		t.Fatalf("Expected 1, was %v", q)
	}

	if q, err := r.AvailableQuantity(i.StartTime, i.EndTime); err != nil {
		fail(t, err)
	} else if q != 0 {
		t.Fatalf("Expected 0, was %v", q)
	}

	if q, err := r.AvailableQuantity(i.EndTime, r.EndTime); err != nil {
		fail(t, err)
	} else if q != 1 {
		t.Fatalf("Expected 1, was %v", q)
	}
}

func TestOverbook(t *testing.T) {
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

	var ob *core.Overbook
	if !errors.As(i.Store(), &ob) {
		t.Fatal()
	}
}
