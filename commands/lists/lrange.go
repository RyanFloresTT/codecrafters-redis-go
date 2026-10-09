package lists

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

func LRange(c resp.Connection, args []string) resp.Value {
	if len(args) != 3 {
		return resp.Error("wrong number of arguments for 'lrange' command")
	}

	listsMu.Lock()
	defer listsMu.Unlock()

	if _, ok := lists[args[0]]; !ok {
		return resp.Array{}
	}

	start, err := strconv.Atoi(args[1])
	if err != nil {
		return resp.Error("value is not an integer or out of range")
	}

	end, err := strconv.Atoi(args[2])
	if err != nil {
		return resp.Error("value is not an integer or out of range")
	}
	length := len(lists[args[0]])
	if start < 0 {
		start += length
	}
	if end < 0 {
		end += length
	}
	if start < 0 {
		start = 0
	}
	if end >= length {
		end = length - 1
	}
	if length == 0 || start >= length || start > end {
		return resp.Array{}
	}
	values := lists[args[0]][start : end+1]
	response := make(resp.Array, len(values))
	for index, value := range values {
		response[index] = resp.BulkString(value)
	}
	return response
}
