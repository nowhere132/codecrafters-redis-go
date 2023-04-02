package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	store := NewStorage()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn, store)
	}
}

func handleConnection(conn net.Conn, store *Storage) {
	defer conn.Close()

	for {
		values, err := DecodeRESP(bufio.NewReader(conn))
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			fmt.Println("error decoding RESP:", err.Error())
			return
		}

		if values.t != Array {
			fmt.Printf("expected array, got first byte %c\n", values.t)
			return
		}
		command := values.Array()[0].String()

		switch command {
		case "ping":
			conn.Write([]byte("+PONG\r\n"))
		case "echo":
			args := values.Array()[1:]
			if len(args) < 1 {
				fmt.Printf("message is required for command ECHO")
				return
			}
			if len(args) < 1 || args[0].t != BulkString {
				fmt.Printf("expected bulk string, got %v\n", args)
				return
			}
			conn.Write([]byte(args[0].EncodeRESP()))
		case "get":
			args := values.Array()[1:]
			if len(args) < 1 {
				fmt.Printf("key is required for command GET")
				return
			}
			key := args[0].String()

			val := store.Get(key)
			conn.Write([]byte(fmt.Sprintf("+%s\r\n", val)))
		case "set":
			args := values.Array()[1:]
			if len(args) < 2 {
				fmt.Printf("key and value are required for command SET")
				return
			}
			key := args[0].String()
			val := args[1].String()

			store.Set(key, val)
			conn.Write([]byte("+OK\r\n"))
		default:
			conn.Write([]byte("-ERR unknown command '" + command + "'\r\n"))
		}
	}
}
