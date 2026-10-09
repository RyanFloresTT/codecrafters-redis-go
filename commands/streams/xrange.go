package streams

import (
	"github.com/codecrafters-io/redis-starter-go/resp"
)

func XRange(c resp.Connection, args []string) resp.Value {
	if len(args) != 3 {
		return resp.Error("wrong number of arguments for 'xrange' command")
	}

	start, _, err := tryParseRangeID(args[1])
	if err != nil {
		return resp.Error(err.Error())
	}
	end, _, err := tryParseRangeID(args[2])
	if err != nil {
		return resp.Error(err.Error())
	}

	response := resp.Array{}
	for _, item := range streamMap[args[0]] {
		if rangeIDBefore(item.id, start) || rangeIDBefore(end, item.id) {
			continue
		}
		response = append(response, item.RESP())
	}
	return response
}
