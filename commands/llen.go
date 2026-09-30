package commands

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

var Llen = Command{
	Name: "LLEN",
	Execute: func(connection net.Conn, args []string) error {
		key := args[0]
		length := len(lists[key])

		err := (&helpers.Connection{Conn: connection}).SendArrayLength(length)

		return err
	},
}
