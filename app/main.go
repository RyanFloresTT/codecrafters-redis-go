package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/cli/info"
	"github.com/codecrafters-io/redis-starter-go/commands"
)

func main() {
	args := os.Args[1:]

	for len(args) > 0 {
		name := strings.ToUpper(strings.TrimPrefix(args[0], "--"))
		command, ok := commands.CLI[name]
		if !ok {
			fmt.Fprintf(os.Stderr, "unknown command %q\n", args[0])
			os.Exit(1)
		}

		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "%s requires a value\n", args[0])
			os.Exit(1)
		}
		_, err := command.Execute(args[1:2])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		args = args[2:]
	}

	port := fmt.Sprintf("%d", info.Redis.Port)

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
