package state

var WatchedKeys = make(map[string]bool)

func ClearWatchedKeys() {
	for key := range WatchedKeys {
		delete(WatchedKeys, key)
	}
}
