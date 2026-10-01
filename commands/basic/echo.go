package basic

import (
	"github.com/codecrafters-io/redis-starter-go/helpers"
)

func Echo(c helpers.Connection, args []string) error {
	return c.SendBulk(args[0])
}
