package transactions

import (
	"github.com/codecrafters-io/redis-starter-go/helpers"
)

func Exec(c helpers.Connection, args []string) error {
	if !ClientTransactions[c.RemoteAddr()] {
		return c.SendError("EXEC without MULTI")
	}

	ClientTransactions[c.RemoteAddr()] = false

	data := make(helpers.Array, len(queued))

	return data.SendTo(c)
}
