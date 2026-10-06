package commands

import (
	"github.com/codecrafters-io/redis-starter-go/cli"
	"github.com/codecrafters-io/redis-starter-go/cli/info"
	"github.com/codecrafters-io/redis-starter-go/commands/basic"
	"github.com/codecrafters-io/redis-starter-go/commands/lists"
	"github.com/codecrafters-io/redis-starter-go/commands/numbers"
	"github.com/codecrafters-io/redis-starter-go/commands/optimistic_locking"
	"github.com/codecrafters-io/redis-starter-go/commands/streams"
	"github.com/codecrafters-io/redis-starter-go/commands/strings"
	"github.com/codecrafters-io/redis-starter-go/commands/transactions"
	"github.com/codecrafters-io/redis-starter-go/commands/types"
)

var Registry = map[string]Command{
	"PING":    {Name: "PING", Execute: basic.Ping},
	"GET":     {Name: "GET", Execute: strings.Get},
	"SET":     {Name: "SET", Execute: strings.Set},
	"ECHO":    {Name: "ECHO", Execute: basic.Echo},
	"RPUSH":   {Name: "RPUSH", Execute: lists.RPush},
	"LRANGE":  {Name: "LRANGE", Execute: lists.LRange},
	"LPUSH":   {Name: "LPUSH", Execute: lists.LPush},
	"LLEN":    {Name: "LLEN", Execute: lists.Llen},
	"LPOP":    {Name: "LPOP", Execute: lists.LPop},
	"BLPOP":   {Name: "BLPOP", Execute: lists.BLPop},
	"TYPE":    {Name: "TYPE", Execute: types.Type},
	"XADD":    {Name: "XADD", Execute: streams.XAdd},
	"XRANGE":  {Name: "XRANGE", Execute: streams.XRange},
	"XREAD":   {Name: "XREAD", Execute: streams.XRead},
	"INCR":    {Name: "INCR", Execute: numbers.Incr},
	"MULTI":   {Name: "MULTI", Execute: transactions.Multi, IsExemptFromQueue: true},
	"EXEC":    {Name: "EXEC", Execute: transactions.Exec, IsExemptFromQueue: true},
	"DISCARD": {Name: "DISCARD", Execute: transactions.Discard, IsExemptFromQueue: true},
	"WATCH":   {Name: "WATCH", Execute: optimistic_locking.Watch, IsExemptFromQueue: true},
	"UNWATCH": {Name: "UNWATCH", Execute: optimistic_locking.Unwatch},
	"INFO":    {Name: "INFO", Execute: info.GetInfo},
}

var CLI = map[string]CLICommand{
	"PORT":      {Name: "PORT", Execute: cli.Port},
	"REPLICAOF": {Name: "REPLICAOF", Execute: cli.ReplicaOf},
}
