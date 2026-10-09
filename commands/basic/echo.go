package basic

import (
	"github.com/codecrafters-io/redis-starter-go/resp"
)

func Echo(c resp.Connection, args []string) resp.Value {
	return resp.BulkString(args[0])
}
