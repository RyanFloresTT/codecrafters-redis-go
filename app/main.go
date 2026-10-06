package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

func main() {
	args := os.Args[1:]

	port := "6379"

	for i := 0; i < len(args); i++ {
		if args[i] == "--port" && i+1 < len(args) {

			// validate port
			if portNum, err := strconv.Atoi(args[i+1]); err != nil || portNum <= 0 || portNum > 65535 {
				fmt.Println("Invalid port:", args[i+1])
				os.Exit(1)
			}

			port = args[i+1]
			i++ // Skip the port value
		}
	}

	l, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		fmt.Println("Failed to bind to port", port)
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
