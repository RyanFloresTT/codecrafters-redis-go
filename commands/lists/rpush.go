package lists

import (
	"sync"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

var lists = make(map[string][]string)
var listsMu sync.Mutex
var listCond = sync.Cond{L: &listsMu}

func RPush(c helpers.Connection, args []string) error {
	key := args[0]
	values := args[1:]

	listsMu.Lock()

	lists[key] = append(lists[key], values...)
	length := len(lists[key])

	listCond.Broadcast()
	listsMu.Unlock()

	return c.SendArrayLength(length)
}
