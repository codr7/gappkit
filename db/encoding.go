package db

import (
	"bufio"
	"io"
	"strconv"
)

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


