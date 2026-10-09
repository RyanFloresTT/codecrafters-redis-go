package transactions

import (
	"github.com/codecrafters-io/redis-starter-go/commands/optimistic_locking/state"
	"github.com/codecrafters-io/redis-starter-go/resp"
)

func Exec(c resp.Connection, args []string) resp.Value {
	for _, wasModified := range state.WatchedKeys {
		if wasModified {
			Discard(c, nil)
			return resp.NullArray{}
		}
	}

	clientsMu.Lock()
	transaction := clients[c.RemoteAddr()]
	delete(clients, c.RemoteAddr())
	clientsMu.Unlock()
	if transaction == nil {
		return resp.Error("EXEC without MULTI")
	}

	data := make(resp.Array, 0, len(transaction.commands))
	for _, command := range transaction.commands {
		response := command.execute(c, command.args)
		data = append(data, response)
	}
	return data
}
