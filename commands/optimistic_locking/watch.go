package optimistic_locking

import (
	"github.com/codecrafters-io/redis-starter-go/commands/optimistic_locking/state"
	"github.com/codecrafters-io/redis-starter-go/commands/transactions"
	"github.com/codecrafters-io/redis-starter-go/resp"
)

// WatchedKeys keeps track of the keys being watched for optimistic locking.
// False by default, indicating that the key has not been modified.
var WatchedKeys = state.WatchedKeys

func Watch(c resp.Connection, args []string) (resp.Value, error) {
	if transactions.IsActive(c) {
		return resp.Error("WATCH inside MULTI is not allowed"), nil
	}
	WatchedKeys[args[0]] = false

	return resp.SimpleString("OK"), nil
}
