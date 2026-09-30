package lists

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

func BLPop(connection net.Conn, args []string) error {
	err := error(nil)
	key := args[0]
	timeout := args[1]

	listsMu.Lock()

	for len(lists[key]) == 0 {
		listCond.Wait()
	}

	fmt.Println("Timeout set to: " + timeout)

	poppedValue := lists[key][0]
	lists[key] = lists[key][1:]

	listsMu.Unlock()

	err = (&helpers.Connection{Conn: connection}).SendArray([]string{key, poppedValue})
	return err
}
