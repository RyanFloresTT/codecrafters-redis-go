package commands

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

var Echo = Command{
	Name: "ECHO",
	Execute: func(connection net.Conn, args []string) error {
		return (helpers.Connection{Conn: connection}).SendBulk(args[0])
	},
}
