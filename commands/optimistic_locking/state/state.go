package state

type WatchedKeysMap map[string]bool

var WatchedKeys = make(WatchedKeysMap)

func (w WatchedKeysMap) Clear() {
	WatchedKeys = make(WatchedKeysMap)
}
