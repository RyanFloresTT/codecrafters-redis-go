package lists

import "github.com/codecrafters-io/redis-starter-go/resp"

func LPush(c resp.Connection, args []string) resp.Value {
	if len(args) < 2 {
		return resp.Error("wrong number of arguments for 'lpush' command")
	}

	key := args[0]
	values := args[1:]

	listsMu.Lock()
	defer listsMu.Unlock()

	for _, value := range values {
		lists[key] = append([]string{value}, lists[key]...)
	}

	listCond.Broadcast()
	return resp.Integer(len(lists[key]))
}
