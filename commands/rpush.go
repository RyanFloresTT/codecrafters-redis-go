package commands

import (
	"net"
	"strconv"
)

var lists = make(map[string][]string)

var RPush = Command{
	Name: "RPUSH",
	Execute: func(connection net.Conn, args []string) error {
		err := error(nil)

		key := args[0]
		values := args[1:]
		lists[key] = append(lists[key], values...)

		_, err = connection.Write([]byte(":" + strconv.Itoa(len(lists[key])) + "\r\n"))

		return err
	},
}
