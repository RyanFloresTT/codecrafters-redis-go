package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/cli/info"
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

func handleMasterConnection(c net.Conn, reader *bufio.Reader) {
	defer c.Close()

	for {
		args, bytesRead, err := parseArgs(reader)
		fmt.Println(bytesRead)
		if err != nil {
			return
		}

		if len(args) == 3 &&
			strings.EqualFold(args[0], "REPLCONF") &&
			strings.EqualFold(args[1], "GETACK") &&
			args[2] == "*" {

			connection := resp.Connection{Conn: c}
			currentOffset := 0 + info.Redis.GetOffset()
			response := resp.Array([]resp.Value{
				resp.BulkString("REPLCONF"),
				resp.BulkString("ACK"),
				resp.BulkString(strconv.Itoa(currentOffset)),
			})
			if err := response.SendTo(connection); err != nil {
				return
			}
			info.Redis.AddToOffset(bytesRead)
			replication.SetReplicaOffset(connection, currentOffset)
			continue
		}

		if _, err := dispatchCommand(resp.Connection{Conn: c}, args, bytesRead); err != nil {
			return
		}
		info.Redis.AddToOffset(bytesRead)
		replication.SetReplicaOffset(resp.Connection{Conn: c}, info.Redis.GetOffset())
	}
}

func readFromConnection(c net.Conn, reader *bufio.Reader) error {
	args, bytesRead, err := parseArgs(reader)
	if err != nil {
		return err
	}

	connection := resp.Connection{Conn: c}
	if len(args) == 0 {
		return connection.SendError("empty command")
	}

	response, err := dispatchCommand(connection, args, bytesRead)
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

func parseArgs(reader *bufio.Reader) ([]string, int, error) {
	header, err := readRESPLine(reader)
	if err != nil {
		return nil, 0, err
	}
	if !strings.HasPrefix(header, "*") {
		return nil, 0, fmt.Errorf("expected RESP array")
	}
	count, err := strconv.Atoi(header[1:])
	if err != nil || count <= 0 {
		return nil, 0, fmt.Errorf("invalid command array length: %s", header)
	}
	args := make([]string, count)
	bytesRead := len(header) + 2
	for index := range args {
		lengthHeader, err := readRESPLine(reader)
		if err != nil {
			return nil, bytesRead, err
		}
		bytesRead += len(lengthHeader) + 2
		if !strings.HasPrefix(lengthHeader, "$") {
			return nil, bytesRead, fmt.Errorf("expected RESP bulk string")
		}
		length, err := strconv.Atoi(lengthHeader[1:])
		if err != nil || length < 0 {
			return nil, bytesRead, fmt.Errorf("invalid bulk string length: %s", lengthHeader)
		}
		data := make([]byte, length)
		if n, err := io.ReadFull(reader, data); err != nil {
			return nil, bytesRead, err
		} else {
			bytesRead += n
		}
		var ending [2]byte
		if n, err := io.ReadFull(reader, ending[:]); err != nil {
			return nil, bytesRead, err
		} else {
			bytesRead += n
		}
		if string(ending[:]) != "\r\n" {
			return nil, bytesRead, fmt.Errorf("invalid bulk string ending")
		}
		args[index] = string(data)
	}
	return args, bytesRead, nil
}

func dispatchCommand(connection resp.Connection, args []string, bytesRead int) (resp.Value, error) {
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
	response := command.Execute(connection, args[1:])

	// kinda don't like this here, as this is called even on replicas, even though it will be empty for them
	for _, replica := range replication.GetReplicas() {
		if !command.Replicates {
			continue
		}
		go func() {
			replicaCommand := command.ToReplicaCommand()
			replicaCommand.SendTo(replica)

			info.Redis.AddToOffset(replicaCommand.ByteLength())
		}()
	}
	return response, nil
}
