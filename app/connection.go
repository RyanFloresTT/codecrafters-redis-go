package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/commands"
	"github.com/codecrafters-io/redis-starter-go/commands/replication"
	"github.com/codecrafters-io/redis-starter-go/commands/transactions"
	"github.com/codecrafters-io/redis-starter-go/resp"
)

func handleConnection(c net.Conn) {
	defer c.Close()
	reader := bufio.NewReader(c)

	for {
		if err := readFromConnection(c, reader); err != nil {
			fmt.Println("Error handling connection: ", err.Error())
			return
		}
	}
}

func handleMasterConnection(c net.Conn) {
	defer c.Close()
	reader := bufio.NewReader(c)

	for {
		args, err := parseArgs(reader)
		if err != nil {
			return
		}

		if _, err := dispatchCommand(resp.Connection{Conn: c}, args); err != nil {
			return
		}
	}
}

func readFromConnection(c net.Conn, reader *bufio.Reader) error {
	args, err := parseArgs(reader)
	if err != nil {
		return err
	}

	connection := resp.Connection{Conn: c}
	if len(args) == 0 {
		return connection.SendError("empty command")
	}

	response, err := dispatchCommand(connection, args)
	if err != nil {
		return err
	}

	return response.SendTo(connection)
}

func readRESPLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	if !strings.HasSuffix(line, "\r\n") {
		return "", fmt.Errorf("invalid RESP line ending")
	}
	return strings.TrimSuffix(line, "\r\n"), nil
}

func parseArgs(reader *bufio.Reader) ([]string, error) {
	header, err := readRESPLine(reader)
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(header, "*") {
		return nil, fmt.Errorf("expected RESP array")
	}
	count, err := strconv.Atoi(header[1:])
	if err != nil || count <= 0 {
		return nil, fmt.Errorf("invalid command array length: %s", header)
	}
	args := make([]string, count)
	for index := range args {
		lengthHeader, err := readRESPLine(reader)
		if err != nil {
			return nil, err
		}
		if !strings.HasPrefix(lengthHeader, "$") {
			return nil, fmt.Errorf("expected RESP bulk string")
		}
		length, err := strconv.Atoi(lengthHeader[1:])
		if err != nil || length < 0 {
			return nil, fmt.Errorf("invalid bulk string length: %s", lengthHeader)
		}
		data := make([]byte, length)
		if _, err := io.ReadFull(reader, data); err != nil {
			return nil, err
		}
		var ending [2]byte
		if _, err := io.ReadFull(reader, ending[:]); err != nil {
			return nil, err
		}
		if string(ending[:]) != "\r\n" {
			return nil, fmt.Errorf("invalid bulk string ending")
		}
		args[index] = string(data)
	}
	return args, nil
}

func dispatchCommand(connection resp.Connection, args []string) (resp.Value, error) {
	name := strings.ToUpper(args[0])
	command, ok := commands.Registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown command '" + args[0] + "'")
	}

	if !command.IsExemptFromQueue && transactions.IsActive(connection) {
		transactions.AddToQueue(connection, command.Execute, args[1:])
		return resp.SimpleString("QUEUED"), nil
	}

	command.Args = args[1:]
	response, err := command.Execute(connection, args[1:])
	if err != nil {
		return nil, err
	}

	for _, replica := range replication.GetReplicas() {
		if !command.Replicates {
			continue
		}
		go func() {
			replicaCommand := command.ToReplicaCommand()
			replicaCommand.SendTo(replica)
		}()
	}
	return response, nil
}
