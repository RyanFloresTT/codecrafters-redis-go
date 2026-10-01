package lists

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

func LPush(c helpers.Connection, args []string) error {
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

	err = c.SendArrayLength(len(lists[key]))

	return err
}
