package db

import (
	"bufio"
	"io"
	"strconv"
	"time"
)

func EncodeBool(val bool, out io.Writer) error {
	var b byte

	if val {
		b = 1
	}

	_, err := out.Write([]byte{b})
	return err
}

func DecodeBool(in *bufio.Reader) (bool, error) {
	v, err := in.ReadByte()

	if v == 1 {
		return true, nil
	}

	return false, err
}

func EncodeInt(val int64, out io.Writer) error {
	v := []byte(strconv.FormatInt(val, 10))
	_, err := out.Write(append([]byte{byte(len(v))}, v...))
	return err
}

func DecodeInt(in *bufio.Reader) (int64, error) {
	n, err := in.ReadByte()
	
	if err != nil {
		return -1, err
	}

	v := make([]byte, n)
	
	if _, err = in.Read(v); err != nil {
		return -1, err
	}

	return strconv.ParseInt(string(v), 10, 64)
}

func EncodeLen(val []byte, out io.Writer) error {
	return EncodeInt(int64(len(val)), out)
}

func DecodeLen(in *bufio.Reader) (int, error) {
	v, err := DecodeInt(in)
	return int(v), err
}

func EncodeRecordId(val RecordId, out io.Writer) error {
	return EncodeInt(val, out)
}

func DecodeRecordId(in *bufio.Reader) (RecordId, error) {
	v, err := DecodeInt(in)
	return RecordId(v), err
}

func EncodeString(val string, out io.Writer) error {
	v := []byte(val)
	
	if err := EncodeLen(v, out); err != nil {
		return err
	}

	_, err := out.Write(v)
	return err
}

func DecodeString(in *bufio.Reader) (string, error) {
	l, err := DecodeLen(in)

	if err != nil {
		return "", err
	}

	v := make([]byte, l)

	if _, err = in.Read(v); err != nil {
		return "", err
	}
	
	return string(v), nil
}

func EncodeTime(val time.Time, out io.Writer) error {
	if err := EncodeInt(val.Unix(), out); err != nil {
		return err
	}

	return nil
}

var nilTime time.Time

func DecodeTime(in *bufio.Reader) (time.Time, error) {
	var s int64
	var err error
	
	if s, err = DecodeInt(in); err != nil {
		return nilTime, err
	}

	return time.Unix(s, 0), nil
}
