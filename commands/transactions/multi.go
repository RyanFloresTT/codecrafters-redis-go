package transactions

import (
	"github.com/codecrafters-io/redis-starter-go/helpers"
)

var IsTransaction bool

func Multi(c helpers.Connection, args []string) error {
	IsTransaction = true
	return c.Send("OK")
}
