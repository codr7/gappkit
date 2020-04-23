package db

import (
	"strings"
	"time"
)

type Order = int

type Compare = func(key, val interface{}) Order

const (
	Lt = Order(-1)
	Eq = Order(0)
	Gt = Order(1)
)

func CompareInt64(x, y int64) Order {
	if x < y {
		return Lt
	}

	if x > y {
		return Gt
	}
	
	return Eq
}

func CompareRecordId(x, y RecordId) Order {
	if x < y {
		return Lt
	}

	if x > y {
		return Gt
	}
	
	return Eq
}

func CompareString(x, y string) Order {
	return Order(strings.Compare(x, y))
}

func CompareTime(x, y interface{}) Order {
	xt := x.(time.Time)
	yt := y.(time.Time)
	
	if xt.Before(yt) {
		return Lt
	}

	if xt.After(yt) {
		return Gt
	}
	
	return Eq
}
