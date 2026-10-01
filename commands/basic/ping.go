package basic

import (
	"github.com/codecrafters-io/redis-starter-go/helpers"
)

func Ping(c helpers.Connection, args []string) error {
	return c.Send("PONG")
}
