package transactions

import (
	"github.com/codecrafters-io/redis-starter-go/commands/optimistic_locking/state"
	"github.com/codecrafters-io/redis-starter-go/resp"
)

func Exec(c resp.Connection, args []string) (resp.Value, error) {
	for _, wasModified := range state.WatchedKeys {
		if wasModified {
			Discard(c, nil)
			state.ClearWatchedKeys()
			return resp.NullArray{}, nil
		}
	}

	clientsMu.Lock()
	transaction := clients[c.RemoteAddr()]
	delete(clients, c.RemoteAddr())
	clientsMu.Unlock()
	if transaction == nil {
		return resp.Error("EXEC without MULTI"), nil
	}

	data := make(resp.Array, 0, len(transaction.commands))
	for _, command := range transaction.commands {
		response, err := command.execute(c, command.args)
		if err != nil {
			return nil, err
		}
		data = append(data, response)
	}
	return data, nil
}
