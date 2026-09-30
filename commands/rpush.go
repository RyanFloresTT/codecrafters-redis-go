package commands

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

var lists = make(map[string][]string)

var RPush = Command{
	Name: "RPUSH",
	Execute: func(connection net.Conn, args []string) error {
		key := args[0]
		values := args[1:]
		lists[key] = append(lists[key], values...)

		return (helpers.Connection{Conn: connection}).SendArrayLength(len(lists[key]))
	},
}
