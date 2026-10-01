package streams

import (
	"github.com/codecrafters-io/redis-starter-go/helpers"
)

type entry struct {
	id     string
	values []kvp
}

type kvp struct {
	key   string
	value string
}

var streamMap = make(map[string][]entry)

func GetMap() map[string][]entry {
	return streamMap
}

func XAdd(c helpers.Connection, args []string) error {
	key := args[0]
	id := args[1]
	values := []kvp{}

	for i := 2; i < len(args); i += 2 {
		values = append(values, kvp{
			key:   args[i],
			value: args[i+1],
		})
	}

	streamMap[key] = append(streamMap[key], entry{
		id:     id,
		values: values,
	})

	return c.SendBulk(id)
}
