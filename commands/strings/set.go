package strings

import (
	"strconv"
	"time"

	"github.com/codecrafters-io/redis-starter-go/commands/optimistic_locking"
	"github.com/codecrafters-io/redis-starter-go/resp"
)

func Set(c resp.Connection, args []string) (resp.Value, error) {
	numArgs := len(args)
	dictionary := GetMap()

	switch numArgs {
	case 2:
		// Handle the case where only the key and value are provided

		setEntry(dictionary, args[0], args[1], time.Time{})
		return resp.SimpleString("OK"), nil

	case 4:
		// Handle the case where key/value and expiration are provided

		switch args[2] {
		case "EX": // Expiration time in seconds
			expireSeconds, err := strconv.Atoi(args[3])
			if err != nil {
				return resp.Error("invalid expire time"), nil
			} else {
				setEntry(dictionary, args[0], args[1], time.Now().Add(time.Duration(expireSeconds)*time.Second))
				return resp.SimpleString("OK"), nil
			}
		case "PX": // Expiration time in milliseconds
			expireMilliseconds, err := strconv.Atoi(args[3])
			if err != nil {
				return resp.Error("invalid expire time"), nil
			} else {
				setEntry(dictionary, args[0], args[1], time.Now().Add(time.Duration(expireMilliseconds)*time.Millisecond))
				return resp.SimpleString("OK"), nil
			}
		default:
			return resp.Error("syntax error"), nil
		}

	default:
		// Handle the default case
		return resp.Error("wrong number of arguments for 'SET' command"), nil
	}
}

func setEntry(dictionary map[string]entry, key string, value string, expiresAt time.Time) {
	dictionary[key] = entry{Value: value, ExpiresAt: expiresAt}
	if _, ok := optimistic_locking.WatchedKeys[key]; ok {
		optimistic_locking.WatchedKeys[key] = true
	}
}
