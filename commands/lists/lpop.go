package lists

import (
	"net"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

func LPop(connection net.Conn, args []string) error {
	err := error(nil)
	key := args[0]

	listsMu.Lock()
	defer listsMu.Unlock()

	if _, exists := lists[key]; !exists {
		err = (&helpers.Connection{Conn: connection}).SendNull()
		return err
	}

	if len(lists[key]) == 0 {
		err = (&helpers.Connection{Conn: connection}).SendNull()
		return err
	}

	if len(args) > 1 {
		amountToPop, _ := strconv.Atoi(args[1])
		poppedValues := []string{}

		for i := 0; i < amountToPop && len(lists[key]) > 0; i++ {
			poppedValue := lists[key][0]
			lists[key] = lists[key][1:]

			poppedValues = append(poppedValues, poppedValue)
		}

		err = (&helpers.Connection{Conn: connection}).SendArray(poppedValues)
		return err
	} else {
		poppedValue := lists[key][0]
		lists[key] = lists[key][1:]

		err = (&helpers.Connection{Conn: connection}).SendBulk(poppedValue)
	}

	return err
}
