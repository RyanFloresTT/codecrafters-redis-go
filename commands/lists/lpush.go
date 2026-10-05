package lists

import "github.com/codecrafters-io/redis-starter-go/helpers"

func LPush(c helpers.Connection, args []string) (helpers.Value, error) {
	if len(args) < 2 {
		return helpers.Error("wrong number of arguments for 'lpush' command"), nil
	}

	key := args[0]
	values := args[1:]

	listsMu.Lock()
	defer listsMu.Unlock()

	for _, value := range values {
		lists[key] = append([]string{value}, lists[key]...)
	}

	listCond.Broadcast()
	return helpers.Integer(len(lists[key])), nil
}
