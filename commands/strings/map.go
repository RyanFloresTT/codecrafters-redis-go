package strings

import "time"

type entry struct {
	Value     string
	ExpiresAt time.Time // zero means no expiry
}

var getMap = map[string]entry{}

func GetMap() map[string]entry {
	return getMap
}

func Entry() entry {
	return entry{}
}

func GetEntry(key string) (entry, bool) {
	value, ok := getMap[key]
	if !ok {
		return entry{}, false
	}
	return value, true
}
