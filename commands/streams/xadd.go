package streams

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

var streamMap = make(map[string][]entry)
var mapMu sync.Mutex
var mapCond = sync.Cond{L: &mapMu}

func GetMap() map[string][]entry {
	return streamMap
}

func XAdd(c helpers.Connection, args []string) (helpers.Value, error) {
	if len(args) < 4 || len(args)%2 != 0 {
		return helpers.Error("wrong number of arguments for 'xadd' command"), nil
	}
	key := args[0]
	idString := args[1]

	newEntry := entry{}

	validationErr := func() *streamError {
		mapMu.Lock()
		defer mapMu.Unlock()

		mapHasEntries := len(streamMap[key]) != 0

		var lastEntry *entry
		if mapHasEntries {
			lastEntry = &streamMap[key][len(streamMap[key])-1]
		}

		if idString == "*" {
			GenerateNewId(&newEntry.id, lastEntry)
		} else if err := ValidateAndParseId(&newEntry, idString, mapHasEntries, lastEntry); err != nil {
			return err
		}

		AddEntry(&newEntry.id, key, args)
		return nil
	}()

	if err := validationErr; err != nil {
		return helpers.Error(err.Error()), nil
	}

	return helpers.BulkString(newEntry.id.String()), nil
}

func GenerateNewId(id *id, lastEntry *entry) {
	id.ms = time.Now().UnixMilli()

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

	mapCond.Broadcast()
}

func ValidateAndParseId(newEntry *entry, idString string, mapHasEntries bool, lastEntry *entry) *streamError {
	idParts := strings.Split(idString, "-")

	if len(idParts) != 2 {
		return &streamError{message: InvalidStreamIDError}
	}

	newEntry.id.ms, _ = strconv.ParseInt(idParts[0], 10, 64)
	seqStr := idParts[1]

	if newEntry.id.ms == 0 && seqStr == "0" {
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
