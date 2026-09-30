package commands

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

var BLPop = Command{
	Name: "BLPOP",
	Execute: func(connection net.Conn, args []string) error {
		err := error(nil)
		key := args[0]
		timeout := args[1]

		if _, exists := lists[key]; !exists {
			err = (&helpers.Connection{Conn: connection}).SendNull()
			return err
		}

		if len(lists[key]) == 0 {
			err = (&helpers.Connection{Conn: connection}).SendNull()
			return err
		}

		for len(lists[key]) == 0 {
		}

		fmt.Println("Timeout set to: " + timeout)

		poppedValue := lists[key][0]
		lists[key] = lists[key][1:]

		err = (&helpers.Connection{Conn: connection}).SendBulk(poppedValue)
		return err
	},
}
