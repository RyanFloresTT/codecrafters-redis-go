package commands

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

var LPush = Command{
	Name: "LPUSH",
	Execute: func(connection net.Conn, args []string) error {
		err := error(nil)

		key := args[0]
		values := args[1:]

		for _, value := range values {
			lists[key] = append([]string{value}, lists[key]...)
		}

		fmt.Println(lists[key])

		err = (&helpers.Connection{Conn: connection}).SendArrayLength(len(lists[key]))

		return err
	},
}
