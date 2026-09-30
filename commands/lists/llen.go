package lists

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

func Llen(connection net.Conn, args []string) error {
	key := args[0]

	listsMu.Lock()
	defer listsMu.Unlock()

	length := len(lists[key])

	err := (&helpers.Connection{Conn: connection}).SendArrayLength(length)

	return err
}
