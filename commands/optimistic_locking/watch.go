package optimistic_locking

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/commands/transactions"
	"github.com/codecrafters-io/redis-starter-go/helpers"
)

// WatchedKeys keeps track of the keys being watched for optimistic locking.
// False by default, indicating that the key has not been modified.
var WatchedKeys = make(map[string]bool)

func Watch(c helpers.Connection, args []string) (helpers.Value, error) {
	fmt.Println("Watching key:", args[0])

	if transactions.IsActive(c) {
		return helpers.Error("WATCH inside MULTI is not allowed"), nil
	}
	WatchedKeys[args[0]] = false

	return helpers.SimpleString("OK"), nil
}
