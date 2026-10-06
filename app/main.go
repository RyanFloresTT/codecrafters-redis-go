package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/cli/info"
	"github.com/codecrafters-io/redis-starter-go/commands"
	"github.com/codecrafters-io/redis-starter-go/resp"
)

func main() {
	args := os.Args[1:]

	parseCLIArgs(&args)

	l := bindToPort()
	checkReplicationRole()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(c)
	}
}

func handleHandshake(connection resp.Connection, port string) {
	handshakeCommands := []resp.Array{
		resp.Array([]resp.Value{resp.BulkString("PING")}),
		resp.Array([]resp.Value{resp.BulkString("REPLCONF"), resp.BulkString("listening-port"), resp.BulkString(port)}),
		resp.Array([]resp.Value{resp.BulkString("REPLCONF"), resp.BulkString("capa"), resp.BulkString("psync2")}),
		resp.Array([]resp.Value{resp.BulkString("PSYNC"), resp.BulkString("?"), resp.BulkString("-1")}),
	}

	for _, cmd := range handshakeCommands {
		cmd.SendTo(connection)

		minBufSize := 1024

		buf := make([]byte, minBufSize)

		_, err := connection.Read(buf)
		if err != nil {
			fmt.Println("Error reading from connection: ", err.Error())
			os.Exit(1)
		}
	}
}

func parseCLIArgs(args *[]string) {
	for len(*args) > 0 {
		name := strings.ToUpper(strings.TrimPrefix((*args)[0], "--"))
		command, ok := commands.CLI[name]
		if !ok {
			fmt.Fprintf(os.Stderr, "unknown command %q\n", (*args)[0])
			os.Exit(1)
		}

		if len(*args) < 2 {
			fmt.Fprintf(os.Stderr, "%s requires a value\n", (*args)[0])
			os.Exit(1)
		}
		err := command.Execute((*args)[1:2])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		*args = (*args)[2:]
	}
}

// Handshake if slave
func checkReplicationRole() {
	if info.Redis.Replication.Role == "slave" {
		address := net.JoinHostPort(info.Redis.Replication.MasterHost, fmt.Sprintf("%d", info.Redis.Replication.MasterPort))
		masterConnection, err := net.DialTimeout("tcp", address, 5*time.Second)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to connect to master:", err)
			os.Exit(1)
		}

		handleHandshake(resp.Connection{Conn: masterConnection}, fmt.Sprintf("%d", info.Redis.Port))
		go handleMasterConnection(masterConnection)
	}
}

func bindToPort() net.Listener {
	port := fmt.Sprintf("%d", info.Redis.Port)
	l, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		fmt.Println("Failed to bind to port", port)
		os.Exit(1)
	}
	return l
}
