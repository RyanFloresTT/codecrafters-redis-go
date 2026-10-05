package commands

import (
	"github.com/codecrafters-io/redis-starter-go/commands/basic"
	"github.com/codecrafters-io/redis-starter-go/commands/lists"
	"github.com/codecrafters-io/redis-starter-go/commands/streams"
	"github.com/codecrafters-io/redis-starter-go/commands/strings"
	"github.com/codecrafters-io/redis-starter-go/commands/types"
	"github.com/codecrafters-io/redis-starter-go/helpers"
)

type Command struct {
	Name    string
	Execute func(connection helpers.Connection, args []string) error
}

var Registry = map[string]Command{
	"PING":   {Name: "PING", Execute: basic.Ping},
	"GET":    {Name: "GET", Execute: strings.Get},
	"SET":    {Name: "SET", Execute: strings.Set},
	"ECHO":   {Name: "ECHO", Execute: basic.Echo},
	"RPUSH":  {Name: "RPUSH", Execute: lists.RPush},
	"LRANGE": {Name: "LRANGE", Execute: lists.LRange},
	"LPUSH":  {Name: "LPUSH", Execute: lists.LPush},
	"LLEN":   {Name: "LLEN", Execute: lists.Llen},
	"LPOP":   {Name: "LPOP", Execute: lists.LPop},
	"BLPOP":  {Name: "BLPOP", Execute: lists.BLPop},
	"TYPE":   {Name: "TYPE", Execute: types.Type},
	"XADD":   {Name: "XADD", Execute: streams.XAdd},
	"XRANGE": {Name: "XRANGE", Execute: streams.XRange},
	"XREAD":  {Name: "XREAD", Execute: streams.XRead},
}
