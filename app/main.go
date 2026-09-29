package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

var getMap = map[string]string{}

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

// *1\r\n$4\r\nPING\r\n

// 1 (number of arguments in the array)

// 4 (length of the argument "PING")

// PING (the argument itself)

/*

Step 1. Read the first line to determine the number of arguments in the array.

Step 2. Read the subsequent lines to get the length and value of each argument.

Step 3. Process each argument as needed (e.g., respond to PING with PONG).

*/

func handleConnection(c net.Conn) {
	defer c.Close()

	fmt.Println("Accepted a connection")

	for {
		if err := readFromConnection(c); err != nil {
			fmt.Println("Error handling connection: ", err.Error())
			return
		}
	}
}

// *2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n

func readFromConnection(c net.Conn) error {
	const minBufSize = 1024

	buf := make([]byte, minBufSize)

	n, err := c.Read(buf)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		os.Exit(1)
	}

	message := string(buf[:n])
	i := strings.Index(message, "\r\n")
	numArgs := (message[1:i])

	numArgsInt, err := strconv.Atoi(numArgs)
	if err != nil {
		fmt.Println("Error converting number of arguments: ", err.Error())
		os.Exit(1)
	}

	fmt.Println("Number of arguments: ", numArgsInt)

	args := make([]string, numArgsInt)

	for j := 0; j < numArgsInt; j++ {
		i += 2 // move past the previous \r\n

		lengthEnd := strings.Index(message[i:], "\r\n")
		if lengthEnd == -1 {
			return fmt.Errorf("incomplete argument length")
		}

		byteLen := message[i : i+lengthEnd]
		byteLenInt, err := strconv.Atoi(byteLen[1:])
		if err != nil {
			return err
		}

		i += lengthEnd + 2
		args[j] = message[i : i+byteLenInt]
		i += byteLenInt
	}

	fmt.Println("Arguments: ", args)

	// write bulk strings ex. $5\r\napple\r\n
	for j := 0; j < numArgsInt; j++ {
		if args[j] == "PING" {
			c.Write([]byte("+PONG\r\n"))
		}

		if args[j] == "ECHO" && j+1 < numArgsInt {
			c.Write([]byte("$" + strconv.Itoa(len(args[j+1])) + "\r\n" + args[j+1] + "\r\n"))
		}

		if args[j] == "SET" {
			if j+2 < numArgsInt {
				getMap[args[j+1]] = args[j+2]
				c.Write([]byte("+OK\r\n"))
			} else {
				c.Write([]byte("-ERR wrong number of arguments for 'SET' command\r\n"))
			}
		}

		if args[j] == "GET" {
			if j+1 < numArgsInt {
				value, ok := getMap[args[j+1]]
				if ok {
					c.Write([]byte("$" + strconv.Itoa(len(value)) + "\r\n" + value + "\r\n"))
				} else {
					c.Write([]byte("$-1\r\n"))
				}
			} else {
				c.Write([]byte("-ERR wrong number of arguments for 'GET' command\r\n"))
			}
		}
	}

	return err
}
