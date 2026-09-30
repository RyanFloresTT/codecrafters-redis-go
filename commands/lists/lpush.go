package lists

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

func LPush(connection net.Conn, args []string) error {
	err := error(nil)

	key := args[0]
	values := args[1:]

	listsMu.Lock()

	for _, value := range values {
		lists[key] = append([]string{value}, lists[key]...)
	}

	listCond.Broadcast()
	listsMu.Unlock()

	fmt.Println(lists[key])

	err = (&helpers.Connection{Conn: connection}).SendArrayLength(len(lists[key]))

	return err
}
