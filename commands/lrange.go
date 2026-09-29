package commands

import (
	"net"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

var LRange = Command{
	Name: "LRANGE",
	Execute: func(connection net.Conn, args []string) error {
		c := helpers.Connection{Conn: connection}

		start, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}
		end, err := strconv.Atoi(args[2])
		if err != nil {
			return err
		}

		if start < 0 || end > len(lists[args[0]]) {
			return c.SendArray([]string{})
		}

		err = c.SendArray(lists[args[0]][start : end+1])
		return err
	},
}
