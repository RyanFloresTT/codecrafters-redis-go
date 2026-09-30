package lists

import (
	"net"
	"strconv"
	"time"

	"vitess.io/vitess/go/timer"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

func BLPop(connection net.Conn, args []string) error {
	err := error(nil)
	key := args[0]
	timeout, _ := strconv.ParseFloat(args[1], 64)

	if timeout > 0 {
		timeout *= 1000
	}

	t := timer.NewTimer(time.Duration(timeout) * time.Millisecond)
	t.Start(func() {
		(&helpers.Connection{Conn: connection}).SendNullArray()
	})

	listsMu.Lock()

	for len(lists[key]) == 0 {
		listCond.Wait()
	}

	poppedValue := lists[key][0]
	lists[key] = lists[key][1:]

	listsMu.Unlock()

	err = (&helpers.Connection{Conn: connection}).SendArray([]string{key, poppedValue})
	return err
}
