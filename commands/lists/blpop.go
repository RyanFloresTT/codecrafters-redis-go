package lists

import (
	"strconv"
	"time"

	"vitess.io/vitess/go/timer"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

func BLPop(c helpers.Connection, args []string) error {
	err := error(nil)
	key := args[0]
	timeout, _ := strconv.ParseFloat(args[1], 64)

	if timeout > 0 {
		timeout *= 1000
	}

	t := timer.NewTimer(time.Duration(timeout) * time.Millisecond)
	t.Start(func() {
		c.SendNullArray()
	})

	listsMu.Lock()

	for len(lists[key]) == 0 {
		listCond.Wait()
	}

	poppedValue := lists[key][0]
	lists[key] = lists[key][1:]

	listsMu.Unlock()

	err = c.SendArray([]string{key, poppedValue})
	return err
}
