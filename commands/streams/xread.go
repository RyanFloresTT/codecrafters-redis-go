package streams

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

func XRead(c helpers.Connection, args []string) error {
	if len(args) != 3 {
		return c.SendError("wrong number of arguments for 'xread' command")
	}
	if !strings.EqualFold(args[0], "streams") {
		return c.SendError("syntax error")
	}
	key := args[1]

	start, err := parseRangeID(args[2], false)
	if err != nil {
		return c.SendError(err.Error())
	}

	data := helpers.Array{}
	for _, item := range streamMap[key] {
		if rangeIDBefore(start, item.id) {
			data = append(data, item.RESP())
		}
	}
	if len(data) == 0 {
		return c.SendNullArray()
	}

	return (helpers.Array{
		helpers.Array{helpers.BulkString(key), data},
	}).SendTo(c)
}
