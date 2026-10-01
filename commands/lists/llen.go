package lists

import (
	"github.com/codecrafters-io/redis-starter-go/helpers"
)

func Llen(c helpers.Connection, args []string) error {
	key := args[0]

	listsMu.Lock()
	defer listsMu.Unlock()

	length := len(lists[key])

	err := c.SendArrayLength(length)

	return err
}
