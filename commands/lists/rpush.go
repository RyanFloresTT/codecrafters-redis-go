package lists

import (
	"net"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

var lists = make(map[string][]string)
var listsMu sync.Mutex
var listCond = sync.Cond{L: &listsMu}

func RPush(connection net.Conn, args []string) error {
	key := args[0]
	values := args[1:]

	listsMu.Lock()

	lists[key] = append(lists[key], values...)
	length := len(lists[key])

	listCond.Broadcast()
	listsMu.Unlock()

	return (helpers.Connection{Conn: connection}).SendArrayLength(length)
}
