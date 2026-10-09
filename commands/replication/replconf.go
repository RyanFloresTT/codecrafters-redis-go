package replication

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/cli/info"
	"github.com/codecrafters-io/redis-starter-go/resp"
)

func ReplConf(c resp.Connection, args []string) resp.Value {
	if args[0] == "GETACK" && args[1] == "*" {
		response := resp.Array([]resp.Value{
			resp.BulkString("REPLCONF "),
			resp.BulkString("ACK"),
			resp.BulkString(strconv.Itoa(info.Redis.GetOffset())),
		})

		return response
	}

	if args[0] == "ACK" {
		replicaOffset, err := strconv.Atoi(args[1])
		if err != nil {
			return resp.SimpleString("ERR invalid ACK offset")
		}
		SetReplicaOffset(c, replicaOffset)
		return resp.Nothing()
	}
	return resp.SimpleString("OK")
}
