package streams

import (
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

var streamMap = make(map[string][]entry)

func GetMap() map[string][]entry {
	return streamMap
}

func XAdd(c helpers.Connection, args []string) error {
	key := args[0]
	idString := args[1]
	values := []kvp{}
	id := id{}

	idParts := strings.Split(idString, "-")
	if len(idParts) != 2 {
		return c.SendError("invalid stream ID")
	}

	mapHasEntries := len(streamMap[key]) != 0
	id.ms = idParts[0]
	seqStr := idParts[1]

	lastEntry := entry{}
	if mapHasEntries {
		lastEntry = streamMap[key][len(streamMap[key])-1]
	}

	if id.ms == "0" && seqStr == "0" {
		return c.SendError("The ID specified in XADD must be greater than 0-0")
	}

	if seqStr != "*" {
		id.seq, _ = strconv.Atoi(idParts[1])
	} else {
		if mapHasEntries {
			if lastEntry.id.ms == id.ms {
				id.seq = lastEntry.id.seq + 1
			} else {
				id.seq = 0
			}
		} else {
			id.seq = 1
		}
	}

	if mapHasEntries {
		if id.ms < lastEntry.id.ms || (id.ms == lastEntry.id.ms && id.seq <= lastEntry.id.seq) {
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

	return c.SendBulk(id.String())
}
