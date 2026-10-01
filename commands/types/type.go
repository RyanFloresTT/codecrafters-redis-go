package types

import (
	"github.com/codecrafters-io/redis-starter-go/commands/streams"
	"github.com/codecrafters-io/redis-starter-go/commands/strings"
	"github.com/codecrafters-io/redis-starter-go/helpers"
)

type Types string

func (t Types) String() string {
	return string(t)
}

const (
	String    Types = "string"
	List      Types = "list"
	Set       Types = "set"
	ZSet      Types = "zset"
	Hash      Types = "hash"
	Stream    Types = "stream"
	VectorSet Types = "vectorset"
	None      Types = "none"
)

func Type(c helpers.Connection, args []string) error {
	key := args[0]
	valueType := None

	if _, ok := strings.GetMap()[key]; ok {
		valueType = String
	} else if _, ok := streams.GetMap()[key]; ok {
		valueType = Stream
	} else {
		valueType = None
	}

	return c.Send(valueType.String())
}
