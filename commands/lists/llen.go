package lists

import (
	"github.com/codecrafters-io/redis-starter-go/resp"
)

func Llen(c resp.Connection, args []string) (resp.Value, error) {
	key := args[0]

	listsMu.Lock()
	defer listsMu.Unlock()

	length := len(lists[key])

	return resp.Integer(length), nil
}
