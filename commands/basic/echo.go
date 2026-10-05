package basic

import (
	"github.com/codecrafters-io/redis-starter-go/helpers"
)

func Echo(c helpers.Connection, args []string) (helpers.Value, error) {
	return helpers.BulkString(args[0]), nil
}
