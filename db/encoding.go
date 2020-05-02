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
	l, err := in.ReadByte()
	
	if err != nil {
		return -1, err
	}

	v := make([]byte, l)
	
	if _, err = in.Read(v); err != nil {
		return -1, err
	}
	
	return strconv.ParseInt(string(v), 10, 64)
}

func EncodeLen(val []byte, out io.Writer) error {
	return EncodeUInt(uint64(len(val)), out)
}

func DecodeLen(in *bufio.Reader) (int, error) {
	v, err := DecodeUInt(in)
	return int(v), err
}

func EncodeRecordId(val RecordId, out io.Writer) error {
	return EncodeUInt(val, out)
}

func DecodeRecordId(in *bufio.Reader) (RecordId, error) {
	v, err := DecodeUInt(in)
	return RecordId(v), err
}

func EncodeUInt(val uint64, out io.Writer) error {
	v := []byte(strconv.FormatUint(val, 10))
	_, err := out.Write(append([]byte{byte(len(v))}, v...))
	return err
}

func DecodeUInt(in *bufio.Reader) (uint64, error) {
	l, err := in.ReadByte()
	
	if err != nil {
		return 0, err
	}

	v := make([]byte, l)
	
	if _, err = in.Read(v); err != nil {
		return 0, err
	}
	
	return strconv.ParseUint(string(v), 10, 64)
}
