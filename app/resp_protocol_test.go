package main

import (
	"bufio"
	"bytes"
	"testing"
)

func TestSimpleString(t *testing.T) {
	val, err := DecodeRESP(readerFrom("+OK\r\n"))

	if err != nil {
		t.Errorf("decode simple string error: %s", err)
	}

	if val.t != SimpleString {
		t.Errorf("expected simple string, got %c", val.t)
	}
	if val.String() != "OK" {
		t.Errorf("expected 'OK', got %s", val.String())
	}
}

func TestBulkString(t *testing.T) {
	val, err := DecodeRESP(readerFrom("$5\r\nhello\r\n"))

	if err != nil {
		t.Errorf("decode bulk string error: %s", err)
	}

	if val.t != BulkString {
		t.Errorf("expected bulk string, got %c", val.t)
	}
	if val.String() != "hello" {
		t.Errorf("expected 'hello', got %s", val.String())
	}
}

func TestDecodeBulkStringArray(t *testing.T) {
	val, err := DecodeRESP(readerFrom("*2\r\n$5\r\nhello\r\n$4\r\nmanh\r\n"))

	if err != nil {
		t.Errorf("decode array error: %s", err)
	}

	if val.t != Array {
		t.Errorf("expected array, got %c", val.t)
	}
	if val.Array()[0].String() != "hello" {
		t.Errorf("expected 'hello', got %s", val.Array()[0].String())
	}
	if val.Array()[1].String() != "manh" {
		t.Errorf("expected 'manh', got %s", val.Array()[1].String())
	}
}

func readerFrom(s string) *bufio.Reader {
	return bufio.NewReader(bytes.NewBufferString(s))
}
