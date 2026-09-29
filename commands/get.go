package commands

import (
	"net"
	"strconv"
	"time"
)

var Get = Command{
	Name: "GET",
	Execute: func(connection net.Conn, args []string) error {
		numArgs := len(args)
		dictionary := GetMap()
		err := error(nil)

		if numArgs > 0 {
			value, ok := dictionary[args[0]]
			if ok {
				if !value.expiresAt.IsZero() && time.Now().After(value.expiresAt) {
					delete(dictionary, args[0])
					ok = false
				}
				if ok {
					connection.Write([]byte("$" + strconv.Itoa(len(value.value)) + "\r\n" + value.value + "\r\n"))
				} else {
					connection.Write([]byte("$-1\r\n"))
				}
			} else {
				connection.Write([]byte("$-1\r\n"))
			}
		} else {
			connection.Write([]byte("-ERR wrong number of arguments for 'GET' command\r\n"))
		}

		return err
	},
}
