package streams

import (
	"math"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

func XRange(c helpers.Connection, args []string) error {
	if len(args) != 3 {
		return c.SendError("wrong number of arguments for 'xrange' command")
	}

	start, err := parseRangeID(args[1], false)
	if err != nil {
		return c.SendError(err.Error())
	}
	end, err := parseRangeID(args[2], true)
	if err != nil {
		return c.SendError(err.Error())
	}

	response := helpers.Array{}
	for _, item := range streamMap[args[0]] {
		if rangeIDBefore(item.id, start) || rangeIDBefore(end, item.id) {
			continue
		}
		response = append(response, item.RESP())
	}
	return response.SendTo(c)
}

func parseRangeID(value string, upperBound bool) (id, error) {
	if value == "-" {
		return id{0, 0}, nil
	}

	if value == "+" {
		return id{ms: math.MaxInt64, seq: math.MaxInt}, nil
	}

	msPart, seqPart, hasSequence := strings.Cut(value, "-")
	ms, err := strconv.ParseInt(msPart, 10, 64)
	if err != nil || ms < 0 {
		return id{}, &streamError{InvalidStreamIDError}
	}

	sequence := 0
	if hasSequence {
		sequence, err = strconv.Atoi(seqPart)
		if err != nil || sequence < 0 {
			return id{}, &streamError{InvalidStreamIDError}
		}
	} else if upperBound {
		sequence = math.MaxInt
	}
	return id{ms, sequence}, nil
}

func rangeIDBefore(left, right id) bool {
	return left.ms < right.ms || (left.ms == right.ms && left.seq < right.seq)
}
