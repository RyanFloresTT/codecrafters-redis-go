package strings

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

func Get(c helpers.Connection, args []string) error {
	numArgs := len(args)
	dictionary := GetMap()

	if numArgs > 0 {
		value, ok := dictionary[args[0]]
		if ok {
			if !value.expiresAt.IsZero() && time.Now().After(value.expiresAt) {
				delete(dictionary, args[0])
				ok = false
			}
			if ok {
				return c.SendBulk(value.value)
			} else {
				return c.SendNull()
			}
		} else {
			return c.SendNull()
		}
	} else {
		return c.SendError("ERR wrong number of arguments for 'GET' command")
	}
}
