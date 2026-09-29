package commands

import (
	"net"
	"strconv"
)

var Echo = Command{
	Name: "ECHO",
	Execute: func(connection net.Conn, args []string) error {
		_, err := connection.Write([]byte("$" + strconv.Itoa(len(args[0])) + "\r\n" + args[0] + "\r\n"))
		return err
	},
}
