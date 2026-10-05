package strings

import (
	"strconv"
	"time"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

func Set(c helpers.Connection, args []string) (helpers.Value, error) {
	numArgs := len(args)
	dictionary := GetMap()

	switch numArgs {
	case 2:
		// Handle the case where only the key and value are provided

		dictionary[args[0]] = entry{Value: args[1]}
		return helpers.SimpleString("OK"), nil

	case 4:
		// Handle the case where key/value and expiration are provided

		switch args[2] {
		case "EX": // Expiration time in seconds
			expireSeconds, err := strconv.Atoi(args[3])
			if err != nil {
				return helpers.Error("invalid expire time"), nil
			} else {
				dictionary[args[0]] = entry{Value: args[1], ExpiresAt: time.Now().Add(time.Duration(expireSeconds) * time.Second)}
				return helpers.SimpleString("OK"), nil
			}
		case "PX": // Expiration time in milliseconds
			expireMilliseconds, err := strconv.Atoi(args[3])
			if err != nil {
				return helpers.Error("invalid expire time"), nil
			} else {
				dictionary[args[0]] = entry{Value: args[1], ExpiresAt: time.Now().Add(time.Duration(expireMilliseconds) * time.Millisecond)}
				return helpers.SimpleString("OK"), nil
			}
		default:
			return helpers.Error("syntax error"), nil
		}

	default:
		// Handle the default case
		return helpers.Error("wrong number of arguments for 'SET' command"), nil
	}
}
