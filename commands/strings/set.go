package strings

import (
	"strconv"
	"time"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

func Set(c helpers.Connection, args []string) error {
	numArgs := len(args)
	dictionary := GetMap()

	switch numArgs {
	case 2:
		// Handle the case where only the key and value are provided

		dictionary[args[0]] = entry{value: args[1]}
		return c.Send("OK")

	case 4:
		// Handle the case where key/value and expiration are provided

		switch args[2] {
		case "EX": // Expiration time in seconds
			expireSeconds, err := strconv.Atoi(args[3])
			if err != nil {
				return c.SendError("ERR invalid expire time")
			} else {
				dictionary[args[0]] = entry{value: args[1], expiresAt: time.Now().Add(time.Duration(expireSeconds) * time.Second)}
				return c.Send("OK")
			}
		case "PX": // Expiration time in milliseconds
			expireMilliseconds, err := strconv.Atoi(args[3])
			if err != nil {
				return c.SendError("ERR invalid expire time")
			} else {
				dictionary[args[0]] = entry{value: args[1], expiresAt: time.Now().Add(time.Duration(expireMilliseconds) * time.Millisecond)}
				return c.Send("OK")
			}
		default:
			return c.SendError("ERR syntax error")
		}

	default:
		// Handle the default case
		return c.SendError("ERR wrong number of arguments for 'SET' command")
	}
}
