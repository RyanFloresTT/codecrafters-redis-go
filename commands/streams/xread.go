package streams

import (
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

func XRead(c helpers.Connection, args []string) error {
	if len(args) < 3 || (len(args)-1)%2 != 0 {
		return c.SendError("wrong number of arguments for 'xread' command")
	}

	timeToBlockMillis := int64(0)
	blocking := false

	if strings.ToUpper(args[0]) == "BLOCK" {
		blocking = true
		timeToBlockMillis, _ = strconv.ParseInt(args[1], 10, 64)

		args = args[2:]
	}

	numOfStreams := (len(args) - 1) / 2

	keys := args[1 : 1+numOfStreams]
	ids := args[1+numOfStreams:]
	response := helpers.Array{}
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

	if ids[0] == "$" {
		lastID, ok := GetLastID(keys[0])
		if ok {
			ids[0] = lastID.String()
		} else {
			ids[0] = "0-0"
		}
	}

	getDataErr := func() *streamError {
		mapMu.Lock()
		defer mapMu.Unlock()

		for {
			response = collectMatches(keys, ids)

			if len(response) > 0 || !blocking || timeout {
				if timer != nil {
					timer.Stop()
				}
				break
			}

			mapCond.Wait()
		}

		return nil
	}()

	if err := getDataErr; err != nil {
		return c.SendError(err.Error())
	}

	if len(response) == 0 {
		return c.SendNullArray()
	}

	if timeout {
		mapCond.Broadcast()
		return c.SendNullArray()
	}

	return response.SendTo(c)
}

func collectMatches(keys []string, ids []string) helpers.Array {
	response := helpers.Array{}

	for index, key := range keys {
		start, _, err := tryParseRangeID(ids[index])
		if err != nil {
			continue
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

	return response
}
