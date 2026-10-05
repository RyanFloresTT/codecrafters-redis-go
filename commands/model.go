package commands

import "github.com/codecrafters-io/redis-starter-go/helpers"

type Command struct {
	Name    string
	Execute func(connection helpers.Connection, args []string) (helpers.Value, error)
	Args    []string
}
