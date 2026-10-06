package replication

import "github.com/codecrafters-io/redis-starter-go/resp"

func ReplConf(c resp.Connection, args []string) (resp.Value, error) {
	if args[0] == "GETACK" && args[1] == "*" {
		response := resp.Array([]resp.Value{
			resp.BulkString("REPLCONF "),
			resp.BulkString("ACK"),
			resp.BulkString("0"),
		})

		return response, nil
	}
	return resp.SimpleString("OK"), nil
}
