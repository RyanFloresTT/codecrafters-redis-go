package types

import (
	"github.com/codecrafters-io/redis-starter-go/commands/streams"
	"github.com/codecrafters-io/redis-starter-go/commands/strings"
	"github.com/codecrafters-io/redis-starter-go/resp"
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

func Type(c resp.Connection, args []string) (resp.Value, error) {
	if len(args) != 1 {
		return resp.Error("wrong number of arguments for 'type' command"), nil
	}
	key := args[0]
	valueType := None

	if _, ok := strings.GetMap()[key]; ok {
		valueType = String
	} else if _, ok := streams.GetMap()[key]; ok {
		valueType = Stream
	} else {
		valueType = None
	}

	return resp.SimpleString(valueType.String()), nil
}
