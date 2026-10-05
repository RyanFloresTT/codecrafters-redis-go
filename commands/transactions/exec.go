package transactions

import (
	"github.com/codecrafters-io/redis-starter-go/helpers"
)

func Exec(c helpers.Connection, args []string) error {
	if !IsTransaction {
		return c.SendError("EXEC without MULTI")
	}

	IsTransaction = false

	data := make(helpers.Array, len(queued))

	return data.SendTo(c)
}
