package optimistic_locking

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/commands/optimistic_locking/state"
	"github.com/codecrafters-io/redis-starter-go/commands/transactions"
	"github.com/codecrafters-io/redis-starter-go/resp"
)

// WatchedKeys keeps track of the keys being watched for optimistic locking.
// False by default, indicating that the key has not been modified.
var WatchedKeys = state.WatchedKeys

func Watch(c resp.Connection, args []string) resp.Value {
	if transactions.IsActive(c) {
		return resp.Error("WATCH inside MULTI is not allowed")
	}

	for _, key := range args {
		WatchedKeys[key] = false
	}

	return resp.SimpleString("OK")
}

func Unwatch(c resp.Connection, args []string) resp.Value {
	fmt.Println(len(WatchedKeys))

	for key := range WatchedKeys {
		delete(WatchedKeys, key)
	}
	return resp.SimpleString("OK")
}
