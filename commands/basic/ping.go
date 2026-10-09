package basic

import (
	"github.com/codecrafters-io/redis-starter-go/resp"
)

func Ping(c resp.Connection, args []string) resp.Value {
	return resp.SimpleString("PONG")
}
