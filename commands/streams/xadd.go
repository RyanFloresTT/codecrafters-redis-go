package streams

import (
	"strings"

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
	// id validation
	idParts := strings.Split(id, "-")
	if len(idParts) != 2 {
		return c.SendError("invalid stream ID")
	}

	ms := idParts[0]
	seq := idParts[1]

	if ms == "0" && seq == "0" {
		return c.SendError("The ID specified in XADD must be greater than 0-0")
	}

	if len(streamMap[key]) != 0 {

		lastEntry := streamMap[key][len(streamMap[key])-1]
		lastEntryIDParts := strings.Split(lastEntry.id, "-")
		lastEntryMS := lastEntryIDParts[0]
		lastEntrySeq := lastEntryIDParts[1]

		if ms < lastEntryMS || (ms == lastEntryMS && seq <= lastEntrySeq) {
			return c.SendError("The ID specified in XADD is equal or smaller than the target stream top item")
		}
	}

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
