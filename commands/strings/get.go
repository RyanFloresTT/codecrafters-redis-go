package strings

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

func Get(c resp.Connection, args []string) (resp.Value, error) {
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
				return resp.BulkString(value.Value), nil
			} else {
				return resp.NullBulk{}, nil
			}
		} else {
			return resp.NullBulk{}, nil
		}
	} else {
		return resp.Error("wrong number of arguments for 'GET' command"), nil
	}
}
