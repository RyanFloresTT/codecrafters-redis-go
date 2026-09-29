package commands

import "net"

var Ping = Command{
	Name: "PING",
	Execute: func(connection net.Conn, args []string) error {
		_, err := connection.Write([]byte("+PONG\r\n"))
		return err
	},
}
