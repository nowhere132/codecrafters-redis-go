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

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
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
			if len(args) != 1 {
				fmt.Printf("expected bulk string, got %v\n", args)
				return
			}
			bulkString := args[0].String()

			resp := fmt.Sprintf("$%d\r\n%s\r\n", len(bulkString), bulkString)
			conn.Write([]byte(resp))
		default:
			conn.Write([]byte("-ERR unknown command '" + command + "'\r\n"))
		}
	}
}
