package basic

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

func Echo(connection net.Conn, args []string) error {
	return (helpers.Connection{Conn: connection}).SendBulk(args[0])
}
