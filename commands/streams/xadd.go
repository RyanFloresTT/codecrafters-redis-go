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

	lastEntry := streamMap[key][len(streamMap[key])-1]

	lastEntryIDParts := strings.Split(lastEntry.id, "-")
	lastEntryMS := lastEntryIDParts[0]
	lastEntrySeq := lastEntryIDParts[1]

	if ms < lastEntryMS {
		return c.SendError("ID specified in XADD is equal or smaller than the target stream top item")
	}

	if ms == lastEntryMS && seq <= lastEntrySeq {
		return c.SendError("ERR The ID specified in XADD must be greater than " + lastEntry.id + "-" + lastEntrySeq)
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
