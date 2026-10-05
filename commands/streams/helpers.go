package streams

import (
	"math"
	"strconv"
	"strings"
)

func rangeIDBefore(left, right id) bool {
	return left.ms < right.ms || (left.ms == right.ms && left.seq < right.seq)
}

func tryParseRangeID(arg string) (id, bool, error) {
	if arg == "-" {
		return id{0, 0}, true, nil
	}
	if arg == "+" {
		return id{ms: math.MaxInt64, seq: math.MaxInt}, true, nil
	}
	msPart, seqPart, hasSequence := strings.Cut(arg, "-")
	ms, err := strconv.ParseInt(msPart, 10, 64)
	if err != nil || ms < 0 {
		return id{}, false, &streamError{InvalidStreamIDError}
	}
	sequence := 0
	if hasSequence {
		sequence, err = strconv.Atoi(seqPart)
		if err != nil || sequence < 0 {
			return id{}, false, &streamError{InvalidStreamIDError}
		}
	}
	return id{ms, sequence}, true, nil
}
