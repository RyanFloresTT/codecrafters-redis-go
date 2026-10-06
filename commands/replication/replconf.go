package replication

import "github.com/codecrafters-io/redis-starter-go/resp"

func ReplConf(c resp.Connection, args []string) (resp.Value, error) {
	return resp.SimpleString("OK"), nil
}
