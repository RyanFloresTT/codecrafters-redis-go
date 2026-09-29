package main

import (
	"fmt"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	fmt.Println("Starting the server...")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(c)
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()

	fmt.Println("Accepted a connection")

	buf := make([]byte, 1024)

	for {
		go readFromConnection(c, buf)
	}
}

func readFromConnection(c net.Conn, buf []byte) (int, error) {
	n, err := c.Read(buf)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		os.Exit(1)
	}

	fmt.Println("Read data from connection: ", string(buf[:n]))

	c.Write([]byte("+PONG\r\n"))

	return n, err
}
