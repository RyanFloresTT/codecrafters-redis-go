package lists

import (
	"sync"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

var lists = make(map[string][]string)
var listsMu sync.Mutex
var listCond = sync.Cond{L: &listsMu}

func RPush(c helpers.Connection, args []string) (helpers.Value, error) {
	if len(args) < 2 {
		return helpers.Error("wrong number of arguments for 'rpush' command"), nil
	}
	key := args[0]
	values := args[1:]

	listsMu.Lock()
	defer listsMu.Unlock()

	lists[key] = append(lists[key], values...)
	length := len(lists[key])

	listCond.Broadcast()
	return helpers.Integer(length), nil
}
