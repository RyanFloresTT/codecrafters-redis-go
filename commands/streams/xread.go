package streams

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

func XRead(c helpers.Connection, args []string) error {
	if len(args) < 3 || (len(args)-1)%2 != 0 {
		return c.SendError("wrong number of arguments for 'xread' command")
	}
	if !strings.EqualFold(args[0], "streams") {
		return c.SendError("syntax error")
	}

	numOfStreams := (len(args) - 1) / 2
	keys := args[1 : 1+numOfStreams]
	ids := args[1+numOfStreams:]
	response := helpers.Array{}

	for index, key := range keys {
		start, _, err := tryParseRangeID(ids[index])
		if err != nil {
			return c.SendError(err.Error())
		}

		data := helpers.Array{}

		for _, item := range streamMap[key] {
			if rangeIDBefore(start, item.id) {
				data = append(data, item.RESP())
			}
		}

		if len(data) > 0 {
			response = append(response, helpers.Array{helpers.BulkString(key), data})
		}
	}

	if len(response) == 0 {
		return c.SendNullArray()
	}
	return response.SendTo(c)
}
