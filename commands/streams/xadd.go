package streams

import (
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

var streamMap = make(map[string][]entry)

func GetMap() map[string][]entry {
	return streamMap
}

func XAdd(c helpers.Connection, args []string) error {
	key := args[0]
	idString := args[1]

	newEntry := entry{}

	mapHasEntries := len(streamMap[key]) != 0

	var lastEntry *entry
	if mapHasEntries {
		lastEntry = &streamMap[key][len(streamMap[key])-1]
	}

	// Case 3 : Auto-Generate Full Id
	if idString == "*" {
		GenerateNewId(&newEntry.id, lastEntry)

		AddEntry(&newEntry.id, key, args)
		return c.SendBulk(newEntry.id.String())
	}

	// Case 1 : Auto-Generate Sequence
	// Case 2 : Validate and Parse Provided Id
	if err := ValidateAndParseId(&newEntry, idString, mapHasEntries, lastEntry); err != nil {
		return c.SendError(err.Error())
	}

	AddEntry(&newEntry.id, key, args)
	return c.SendBulk(newEntry.id.String())
}

func GenerateNewId(id *id, lastEntry *entry) {
	id.ms = strconv.FormatInt(time.Now().UnixMilli(), 10)

	if lastEntry != nil && lastEntry.id.ms == id.ms {
		id.seq = lastEntry.id.seq + 1
	} else {
		id.seq = 0
	}
}

func AddEntry(id *id, key string, args []string) {
	values := []kvp{}

	for i := 2; i < len(args); i += 2 {
		values = append(values, kvp{
			key:   args[i],
			value: args[i+1],
		})
	}

	streamMap[key] = append(streamMap[key], entry{
		id:     *id,
		values: values,
	})
}

func ValidateAndParseId(newEntry *entry, idString string, mapHasEntries bool, lastEntry *entry) *streamError {
	idParts := strings.Split(idString, "-")

	if len(idParts) != 2 {
		return &streamError{message: InvalidStreamIDError}
	}

	newEntry.id.ms = idParts[0]
	seqStr := idParts[1]

	if newEntry.id.ms == "0" && seqStr == "0" {
		return &streamError{message: XAddIDGreaterThanZeroError}
	}

	if seqStr != "*" {
		newEntry.id.seq, _ = strconv.Atoi(idParts[1])
	} else {
		if mapHasEntries {
			if lastEntry.id.ms == newEntry.id.ms {
				newEntry.id.seq = lastEntry.id.seq + 1
			} else {
				newEntry.id.seq = 0
			}
		} else {
			newEntry.id.seq = 1
		}
	}

	if mapHasEntries {
		if newEntry.id.ms < lastEntry.id.ms || (newEntry.id.ms == lastEntry.id.ms && newEntry.id.seq <= lastEntry.id.seq) {
			return &streamError{message: XAddIDSmallerThanTopError}
		}
	}

	return nil
}
