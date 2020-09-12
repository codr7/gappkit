package compare

import (
	"time"
	"unsafe"
)

type Order = int

type Func = func(key, val interface{}) Order

const (
	Lt = Order(-1)
	Eq = Order(0)
	Gt = Order(1)
)

func Bool(x, y bool) Order {
	if !x && y {
		return Lt
	}

	if x && !y {
		return Gt
	}
	
	return Eq
}

func Int(x, y int64) Order {
	if x < y {
		return Lt
	}

	if x > y {
		return Gt
	}
	
	return Eq
}

func Pointer(x, y unsafe.Pointer) Order {
	xp, yp := uintptr(x), uintptr(y)
	
	if xp < yp {
		return Lt
	}
	
	if xp > yp {
		return Gt
	}
	
	return Eq
}

func String(x, y string) Order {
	if x < y {
		return Lt
	}

	if x > y {
		return Gt
	}
	
	return Eq
}

func Time(x, y interface{}) Order {
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
