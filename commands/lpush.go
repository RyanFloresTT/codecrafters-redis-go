package commands

import (
	"fmt"
	"net"
	"strconv"
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

		_, err = connection.Write([]byte(":" + strconv.Itoa(len(lists[key])) + "\r\n"))

		return err
	},
}
