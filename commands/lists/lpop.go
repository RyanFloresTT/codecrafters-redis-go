package lists

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

func LPop(c resp.Connection, args []string) (resp.Value, error) {
	if len(args) == 0 || len(args) > 2 {
		return resp.Error("wrong number of arguments for 'lpop' command"), nil
	}
	key := args[0]

	listsMu.Lock()
	defer listsMu.Unlock()

	if _, exists := lists[key]; !exists {
		return resp.NullBulk{}, nil
	}

	if len(lists[key]) == 0 {
		return resp.NullBulk{}, nil
	}

	if len(args) > 1 {
		amountToPop, err := strconv.Atoi(args[1])
		if err != nil || amountToPop < 0 {
			return resp.Error("value is out of range, must be positive"), nil
		}
		poppedValues := []string{}

		for i := 0; i < amountToPop && len(lists[key]) > 0; i++ {
			poppedValue := lists[key][0]
			lists[key] = lists[key][1:]

			poppedValues = append(poppedValues, poppedValue)
		}

		response := make(resp.Array, len(poppedValues))
		for index, value := range poppedValues {
			response[index] = resp.BulkString(value)
		}
		return response, nil
	} else {
		poppedValue := lists[key][0]
		lists[key] = lists[key][1:]

		return resp.BulkString(poppedValue), nil
	}
}
