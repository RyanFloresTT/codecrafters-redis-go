package strings

import "time"

type entry struct {
	value     string
	expiresAt time.Time // zero means no expiry
}

var getMap = map[string]entry{}

func GetMap() map[string]entry {
	return getMap
}

func Entry() entry {
	return entry{}
}
