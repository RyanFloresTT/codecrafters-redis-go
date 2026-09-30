package commands

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

var LPop = Command{
	Name: "LPOP",
	Execute: func(connection net.Conn, args []string) error {
		err := error(nil)

		key := args[0]
		if _, exists := lists[key]; !exists {
			err = (&helpers.Connection{Conn: connection}).SendNull()
		}

		if len(lists[key]) > 0 {
			poppedValue := lists[key][0]
			lists[key] = lists[key][1:]
			err = (&helpers.Connection{Conn: connection}).SendBulk(poppedValue)
		} else {
			err = (&helpers.Connection{Conn: connection}).SendNull()
		}

		return err
	},
}
