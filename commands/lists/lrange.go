package lists

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

func LRange(c helpers.Connection, args []string) error {

	listsMu.Lock()
	defer listsMu.Unlock()

	if _, ok := lists[args[0]]; !ok {
		return c.SendArray([]string{})
	}

	start, err := strconv.Atoi(args[1])
	if err != nil {
		return err
	}

	end, err := strconv.Atoi(args[2])
	if err != nil {
		return err
	}

	if end > len(lists[args[0]])-1 {
		end = len(lists[args[0]]) - 1
	}

	if start < 0 {
		start += len(lists[args[0]])
		if start < 0 {
			start = 0
		}
	}

	if end < 0 {
		end += len(lists[args[0]])
		if end < 0 {
			end = 0
		}
	}

	err = c.SendArray(lists[args[0]][start : end+1])
	return err
}
