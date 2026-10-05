package streams

import (
	"github.com/codecrafters-io/redis-starter-go/helpers"
)

func XRange(c helpers.Connection, args []string) (helpers.Value, error) {
	if len(args) != 3 {
		return helpers.Error("wrong number of arguments for 'xrange' command"), nil
	}

	start, _, err := tryParseRangeID(args[1])
	if err != nil {
		return helpers.Error(err.Error()), nil
	}
	end, _, err := tryParseRangeID(args[2])
	if err != nil {
		return helpers.Error(err.Error()), nil
	}

	response := helpers.Array{}
	for _, item := range streamMap[args[0]] {
		if rangeIDBefore(item.id, start) || rangeIDBefore(end, item.id) {
			continue
		}
		response = append(response, item.RESP())
	}
	return response, nil
}
