package strings

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

func Get(c helpers.Connection, args []string) (helpers.Value, error) {
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
				return helpers.BulkString(value.Value), nil
			} else {
				return helpers.NullBulk{}, nil
			}
		} else {
			return helpers.NullBulk{}, nil
		}
	} else {
		return helpers.Error("wrong number of arguments for 'GET' command"), nil
	}
}
