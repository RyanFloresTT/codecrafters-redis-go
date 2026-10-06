package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
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

func handleHandshake(connection resp.Connection, port string, reader *bufio.Reader) error {
	handshakeCommands := []resp.Array{
		resp.Array([]resp.Value{resp.BulkString("PING")}),
		resp.Array([]resp.Value{resp.BulkString("REPLCONF"), resp.BulkString("listening-port"), resp.BulkString(port)}),
		resp.Array([]resp.Value{resp.BulkString("REPLCONF"), resp.BulkString("capa"), resp.BulkString("psync2")}),
		resp.Array([]resp.Value{resp.BulkString("PSYNC"), resp.BulkString("?"), resp.BulkString("-1")}),
	}

	for index, cmd := range handshakeCommands {
		if err := cmd.SendTo(connection); err != nil {
			return err
		}
		line, err := readRESPLine(reader)
		if err != nil {
			return err
		}
		if index == len(handshakeCommands)-1 {
			if !strings.HasPrefix(line, "+FULLRESYNC ") {
				return fmt.Errorf("unexpected PSYNC response: %s", line)
			}
		} else {
			expected := "+OK"
			if index == 0 {
				expected = "+PONG"
			}
			if line != expected {
				return fmt.Errorf("unexpected handshake response: %s", line)
			}
		}
	}
	header, err := readRESPLine(reader)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(header, "$") {
		return fmt.Errorf("expected RDB length, got %s", header)
	}
	length, err := strconv.ParseInt(header[1:], 10, 64)
	if err != nil || length < 0 {
		return fmt.Errorf("invalid RDB length: %s", header)
	}
	_, err = io.CopyN(io.Discard, reader, length)
	return err
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

		reader := bufio.NewReader(masterConnection)
		if err := handleHandshake(resp.Connection{Conn: masterConnection}, fmt.Sprintf("%d", info.Redis.Port), reader); err != nil {
			masterConnection.Close()
			fmt.Fprintln(os.Stderr, "Failed replication handshake:", err)
			os.Exit(1)
		}
		go handleMasterConnection(masterConnection, reader)
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
