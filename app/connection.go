package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/commands"
	"github.com/codecrafters-io/redis-starter-go/helpers"
)

func handleConnection(c net.Conn) {
	defer c.Close()

	for {
		if err := readFromConnection(c); err != nil {
			fmt.Println("Error handling connection: ", err.Error())
			return
		}
	}
}

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

	connection := helpers.Connection{Conn: c}

	err = commands.Registry[args[0]].Execute(connection, args[1:])

	return err
}
