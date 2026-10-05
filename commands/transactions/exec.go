package transactions

import (
	"github.com/codecrafters-io/redis-starter-go/helpers"
)

func Exec(c helpers.Connection, args []string) (helpers.Value, error) {
	clientsMu.Lock()
	transaction := clients[c.RemoteAddr()]
	delete(clients, c.RemoteAddr())
	clientsMu.Unlock()
	if transaction == nil {
		return helpers.Error("EXEC without MULTI"), nil
	}

	data := make(helpers.Array, 0, len(transaction.commands))
	for _, command := range transaction.commands {
		response, err := command.execute(c, command.args)
		if err != nil {
			return nil, err
		}
		data = append(data, response)
	}
	return data, nil
}
