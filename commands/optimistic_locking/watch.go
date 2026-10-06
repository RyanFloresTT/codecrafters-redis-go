package optimistic_locking

import (
	"github.com/codecrafters-io/redis-starter-go/helpers"
)

var WatchedKeys = []string{}

func Watch(c helpers.Connection, args []string) (helpers.Value, error) {
	WatchedKeys = append(WatchedKeys, args[0])

	return helpers.SimpleString("OK"), nil
}
