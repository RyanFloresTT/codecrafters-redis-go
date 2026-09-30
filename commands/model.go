package commands

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/commands/basic"
	"github.com/codecrafters-io/redis-starter-go/commands/lists"
	"github.com/codecrafters-io/redis-starter-go/commands/strings"
)

type Command struct {
	Name    string
	Execute func(connection net.Conn, args []string) error
}

var Ping = Command{Name: "PING", Execute: basic.Ping}
var Echo = Command{Name: "ECHO", Execute: basic.Echo}
var Get = Command{Name: "GET", Execute: strings.Get}
var Set = Command{Name: "SET", Execute: strings.Set}
var GetMap = strings.GetMap
var Entry = strings.Entry
var RPush = Command{Name: "RPUSH", Execute: lists.RPush}
var LPush = Command{Name: "LPUSH", Execute: lists.LPush}
var Llen = Command{Name: "LLEN", Execute: lists.Llen}
var LPop = Command{Name: "LPOP", Execute: lists.LPop}
var LRange = Command{Name: "LRANGE", Execute: lists.LRange}
var BLPop = Command{Name: "BLPOP", Execute: lists.BLPop}

var Registry = map[string]Command{
	"PING":   Ping,
	"GET":    Get,
	"SET":    Set,
	"ECHO":   Echo,
	"RPUSH":  RPush,
	"LRANGE": LRange,
	"LPUSH":  LPush,
	"LLEN":   Llen,
	"LPOP":   LPop,
	"BLPOP":  BLPop,
}
