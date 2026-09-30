package basic

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

func Ping(connection net.Conn, args []string) error {
	return (helpers.Connection{Conn: connection}).Send("PONG")
}
