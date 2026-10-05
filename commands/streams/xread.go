package streams

import (
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

func XRead(c helpers.Connection, args []string) (helpers.Value, error) {
	if len(args) < 3 || (len(args)-1)%2 != 0 {
		return helpers.Error("wrong number of arguments for 'xread' command"), nil
	}

	timeToBlockMillis := int64(0)
	blocking := false

	if strings.ToUpper(args[0]) == "BLOCK" {
		blocking = true
		var err error
		timeToBlockMillis, err = strconv.ParseInt(args[1], 10, 64)
		if err != nil || timeToBlockMillis < 0 {
			return helpers.Error("timeout is not an integer or out of range"), nil
		}

		args = args[2:]
	}
	if len(args) < 3 || !strings.EqualFold(args[0], "STREAMS") || (len(args)-1)%2 != 0 {
		return helpers.Error("syntax error"), nil
	}

	numOfStreams := (len(args) - 1) / 2

	keys := args[1 : 1+numOfStreams]
	ids := args[1+numOfStreams:]
	timeout := false

	var timer *time.Timer

	if timeToBlockMillis > 0 {
		timer = time.AfterFunc(time.Duration(timeToBlockMillis)*time.Millisecond, func() {
			mapMu.Lock()
			timeout = true
			mapCond.Broadcast()
			mapMu.Unlock()
		})
		defer timer.Stop()
	}

	var response helpers.Array
	var readErr error
	func() {
		mapMu.Lock()
		defer mapMu.Unlock()

		for index, key := range keys {
			if ids[index] == "$" {
				if lastID, ok := GetLastID(key); ok {
					ids[index] = lastID.String()
				} else {
					ids[index] = "0-0"
				}
			}
		}

		for {
			response, readErr = collectMatches(keys, ids)
			if readErr != nil {
				return
			}

			if len(response) > 0 || !blocking || timeout {
				return
			}

			mapCond.Wait()
		}
	}()

	if readErr != nil {
		return helpers.Error(readErr.Error()), nil
	}

	if len(response) == 0 {
		return helpers.NullArray{}, nil
	}

	return response, nil
}

func collectMatches(keys []string, ids []string) (helpers.Array, error) {
	response := helpers.Array{}

	for index, key := range keys {
		start, _, err := tryParseRangeID(ids[index])
		if err != nil {
			return nil, err
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

	return response, nil
}
