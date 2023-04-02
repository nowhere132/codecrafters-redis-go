package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type Type byte

const (
	SimpleString Type = '+'
	BulkString   Type = '$'
	Array        Type = '*'
)

type Value struct {
	t        Type
	bytes    []byte
	children []*Value
}

func (v *Value) EncodeRESP() string {
	switch v.t {
	case SimpleString:
		return fmt.Sprintf("+%s\r\n", v.String())
	case BulkString:
		s := v.String()
		return fmt.Sprintf("$%d\r\n%s\r\n", len(s), s)
	case Array:
		children := v.Array()
		s := fmt.Sprintf("*%d\r\n", len(children))
		for _, child := range children {
			s += child.EncodeRESP()
		}
		return s
	default:
		return ""
	}
}

func DecodeRESP(byteStream *bufio.Reader) (*Value, error) {
	typeByte, err := byteStream.ReadByte()
	if err != nil {
		return nil, err
	}

	switch typeByte {
	case '+':
		return decodeSimpleString(byteStream)
	case '$':
		return decodeBulkString(byteStream)
	case '*':
		return decodeArray(byteStream)
	}

	return nil, fmt.Errorf("invalid RESP data type byte: %c", typeByte)
}

func decodeSimpleString(byteStream *bufio.Reader) (*Value, error) {
	readBytes, err := readUntilCRLF(byteStream)
	if err != nil {
		return nil, err
	}

	node := Value{t: SimpleString, bytes: readBytes}
	return &node, nil
}

func decodeBulkString(byteStream *bufio.Reader) (*Value, error) {
	readBytesForLength, err := readUntilCRLF(byteStream)
	if err != nil {
		return nil, fmt.Errorf("failed to read bulk string length: %s", err)
	}

	count, err := strconv.Atoi(string(readBytesForLength))
	if err != nil {
		return nil, fmt.Errorf("failed to parse bulk string length: %s", err)
	}

	readBytes := make([]byte, count+2)
	if _, err := io.ReadFull(byteStream, readBytes); err != nil {
		return nil, err
	}

	node := Value{t: BulkString, bytes: readBytes[:count]}
	return &node, nil
}

func decodeArray(byteStream *bufio.Reader) (*Value, error) {
	readBytesForLength, err := readUntilCRLF(byteStream)
	if err != nil {
		return nil, fmt.Errorf("failed to read array length: %s", err)
	}

	count, err := strconv.Atoi(string(readBytesForLength))
	if err != nil {
		return nil, fmt.Errorf("failed to parse array length: %s", err)
	}

	var arr []*Value
	for i := 0; i < count; i++ {
		val, err := DecodeRESP(byteStream)
		if err != nil {
			return nil, fmt.Errorf("failed to decode element %d in array: %s", i, err)
		}

		arr = append(arr, val)
	}

	node := Value{t: Array, children: arr}
	return &node, nil
}

func readUntilCRLF(byteStream *bufio.Reader) ([]byte, error) {
	var result []byte

	for {
		bytes, err := byteStream.ReadBytes('\n')
		if err != nil {
			return nil, err
		}

		result = append(result, bytes...)
		if len(result) >= 2 && result[len(result)-2] == '\r' {
			break
		}
	}

	if result[len(result)-2] != '\r' {
		return nil, fmt.Errorf("couldn't find CRLR in %s", result)
	}

	return result[:len(result)-2], nil
}

func (v *Value) String() string {
	if v.t == SimpleString || v.t == BulkString {
		return string(v.bytes)
	}
	return ""
}

func (v *Value) Array() []*Value {
	if v.t == Array {
		return v.children
	}
	return []*Value{}
}
