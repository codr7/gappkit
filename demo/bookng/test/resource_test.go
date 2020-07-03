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

func TestResource(t *testing.T) {
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

func TestCategory(t *testing.T) {
	setup(t)
	defer teardown(t)

	foo := db.NewResource()
	foo.Name = "foo"
	foo.Quantity = 0
	
	if err := foo.Store(); err != nil {
		fail(t, err)
	}

	bar := db.NewResource()
	bar.Name = "bar"
	bar.AddCategory(foo.Id())
	bar.AddCategory(foo.Id())
	
	if err := bar.Store(); err != nil {
		fail(t, err)
	}
	
	i := db.NewItem()
	i.Resource = bar.Id()

	if err := i.Store(); err != nil {
		fail(t, err)
	}

	if q, err := foo.AvailableQuantity(foo.StartTime, i.StartTime); err != nil {
		fail(t, err)
	} else if q != 2 {
		t.Fatalf("Expected 2, was %v", q)
	}

	if q, err := foo.AvailableQuantity(i.StartTime, i.EndTime); err != nil {
		fail(t, err)
	} else if q != 0 {
		t.Fatalf("Expected 0, was %v", q)
	}

	if q, err := foo.AvailableQuantity(i.EndTime, foo.EndTime); err != nil {
		fail(t, err)
	} else if q != 2 {
		t.Fatalf("Expected 2, was %v", q)
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
