package rdb

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

func ConfigCommand(c resp.Connection, args []string) resp.Value {
	command := strings.ToUpper(args[0])

	if command == "GET" {
		searchValue := args[1]

		if searchValue == "dir" {
			return resp.CommandArray{
				resp.BulkString("dir"),
				resp.BulkString(RDBConfig.Dir),
			}
		}
		if searchValue == "dbfilename" {
			return resp.CommandArray{
				resp.BulkString("dbfilename"),
				resp.BulkString(RDBConfig.DbFileName),
			}
		}
	}

	if command == "SET" {
		// Handle the SET subcommand
	}

	return resp.Error("unsupported CONFIG command")
}
