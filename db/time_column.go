package db

import (
	"bufio"
	"gappkit/compare"
	"io"
	"time"
)

type TimeColumn struct {
	BasicColumn
}

func (self *TimeColumn) Init(name string) *TimeColumn {
	self.BasicColumn.Init(name)
	return self
}

func (self *TimeColumn) Compare(x, y interface{}) compare.Order {
	return compare.Time(x.(time.Time), y.(time.Time))
}

func (self *TimeColumn) Decode(in *bufio.Reader) (interface{}, error) {
	var s, ns int64
	var err error
	
	if s, err = DecodeInt(in); err != nil {
		return nil, err
	}

	if ns, err = DecodeInt(in); err != nil {
		return nil, err
	}

	return time.Unix(s, ns).UTC(), nil
}

func (self *TimeColumn) Encode(val interface{}, out io.Writer) error {
	v := val.(time.Time)
	
	if err := EncodeInt(v.Unix(), out); err != nil {
		return err
	}

	if err := EncodeInt(v.UnixNano() % 1000000000, out); err != nil {
		return err
	}

	return nil
}
