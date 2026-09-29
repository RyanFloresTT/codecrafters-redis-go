package commands

import (
	"net"
	"strconv"
	"time"
)

var Set = Command{
	Name: "SET",
	Execute: func(connection net.Conn, args []string) error {
		numArgs := len(args)
		dictionary := GetMap()
		err := error(nil)

		switch numArgs {
		case 2:
			// Handle the case where only the key and value are provided

			dictionary[args[0]] = entry{value: args[1]}
			_, err = connection.Write([]byte("+OK\r\n"))

		case 4:
			// Handle the case where key/value and expiration are provided

			switch args[2] {
			case "EX": // Expiration time in seconds
				expireSeconds, err := strconv.Atoi(args[3])
				if err != nil {
					_, err = connection.Write([]byte("-ERR invalid expire time\r\n"))
				} else {
					dictionary[args[0]] = entry{value: args[1], expiresAt: time.Now().Add(time.Duration(expireSeconds) * time.Second)}
					_, err = connection.Write([]byte("+OK\r\n"))
				}
			case "PX": // Expiration time in milliseconds
				expireMilliseconds, err := strconv.Atoi(args[3])
				if err != nil {
					_, err = connection.Write([]byte("-ERR invalid expire time\r\n"))
				} else {
					dictionary[args[0]] = entry{value: args[1], expiresAt: time.Now().Add(time.Duration(expireMilliseconds) * time.Millisecond)}
					_, err = connection.Write([]byte("+OK\r\n"))
				}
			default:
				_, err = connection.Write([]byte("-ERR syntax error\r\n"))
			}

		default:
			// Handle the default case
			_, err = connection.Write([]byte("-ERR wrong number of arguments for 'SET' command\r\n"))
		}

		return err
	},
}
