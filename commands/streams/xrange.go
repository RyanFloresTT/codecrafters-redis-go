package streams

import (
	"github.com/codecrafters-io/redis-starter-go/helpers"
)

func XRange(c helpers.Connection, args []string) error {
	if len(args) != 3 {
		return c.SendError("wrong number of arguments for 'xrange' command")
	}

	start, _, err := tryParseRangeID(args[1])
	if err != nil {
		return c.SendError(err.Error())
	}
	end, _, err := tryParseRangeID(args[2])
	if err != nil {
		return c.SendError(err.Error())
	}

	response := helpers.Array{}
	for _, item := range streamMap[args[0]] {
		if rangeIDBefore(item.id, start) || rangeIDBefore(end, item.id) {
			continue
		}
		response = append(response, item.RESP())
	}
	return response.SendTo(c)
}
