package commands

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

var Ping = Command{
	Name: "PING",
	Execute: func(connection net.Conn, args []string) error {
		return (helpers.Connection{Conn: connection}).Send("PONG")
	},
}
