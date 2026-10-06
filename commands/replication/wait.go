package replication

import (
	"github.com/codecrafters-io/redis-starter-go/resp"
)

func Wait(c resp.Connection, args []string) (resp.Value, error) {
	return resp.Integer(0), nil
}
