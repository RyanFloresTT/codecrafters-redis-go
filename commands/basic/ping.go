package basic

import (
	"github.com/codecrafters-io/redis-starter-go/helpers"
)

func Ping(c helpers.Connection, args []string) (helpers.Value, error) {
	return helpers.SimpleString("PONG"), nil
}
