package commands

import "net"

type Command struct {
	Name    string
	Execute func(connection net.Conn, args []string) error
}

var Registry = map[string]Command{
	"PING":   Ping,
	"GET":    Get,
	"SET":    Set,
	"ECHO":   Echo,
	"RPUSH":  RPush,
	"LRANGE": LRange,
}
