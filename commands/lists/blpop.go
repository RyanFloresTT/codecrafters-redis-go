package lists

import (
	"strconv"
	"time"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

func BLPop(c resp.Connection, args []string) (resp.Value, error) {
	if len(args) != 2 {
		return resp.Error("wrong number of arguments for 'blpop' command"), nil
	}
	key := args[0]
	timeoutSeconds, err := strconv.ParseFloat(args[1], 64)
	if err != nil || timeoutSeconds < 0 {
		return resp.Error("timeout is not a float or out of range"), nil
	}

	timedOut := false
	var timer *time.Timer
	if timeoutSeconds > 0 {
		timer = time.AfterFunc(time.Duration(timeoutSeconds*float64(time.Second)), func() {
			listsMu.Lock()
			timedOut = true
			listCond.Broadcast()
			listsMu.Unlock()
		})
		defer timer.Stop()
	}

	listsMu.Lock()
	defer listsMu.Unlock()

	for len(lists[key]) == 0 && !timedOut {
		listCond.Wait()
	}
	if len(lists[key]) == 0 {
		return resp.NullArray{}, nil
	}

	poppedValue := lists[key][0]
	lists[key] = lists[key][1:]

	return resp.Array{resp.BulkString(key), resp.BulkString(poppedValue)}, nil
}
