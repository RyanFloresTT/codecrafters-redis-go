package commands

import "github.com/codecrafters-io/redis-starter-go/resp"

type Command struct {
	Name              string
	Execute           func(connection resp.Connection, args []string) (resp.Value, error)
	Args              []string
	IsExemptFromQueue bool
}
