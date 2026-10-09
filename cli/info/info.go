package info

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

func GetInfo(c resp.Connection, args []string) resp.Value {
	if len(args) == 0 || strings.EqualFold(args[0], "replication") {
		return resp.BulkString(Redis.Replication.String())
	}
	return resp.BulkString("")
}
