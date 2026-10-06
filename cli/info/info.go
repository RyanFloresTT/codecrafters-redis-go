package info

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

func GetInfo(c resp.Connection, args []string) (resp.Value, error) {
	if len(args) == 0 || strings.EqualFold(args[0], "replication") {
		return resp.BulkString(Redis.Replication.String()), nil
	}
	return resp.BulkString(""), nil
}
