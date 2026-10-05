package strings

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

func Get(c helpers.Connection, args []string) error {
	numArgs := len(args)
	dictionary := GetMap()

	if numArgs > 0 {
		value, ok := GetEntry(args[0])
		if ok {
			if !value.ExpiresAt.IsZero() && time.Now().After(value.ExpiresAt) {
				delete(dictionary, args[0])
				ok = false
			}
			if ok {
				return c.SendBulk(value.Value)
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
