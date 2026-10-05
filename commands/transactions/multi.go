package transactions

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

var ClientTransactions = make(map[net.Addr]bool)

type CommandFunc func(helpers.Connection, []string) error

type queuedCommand struct {
	execute CommandFunc
	args    []string
}

var queued []queuedCommand

func Multi(c helpers.Connection, args []string) error {
	ClientTransactions[c.Conn.RemoteAddr()] = true
	return c.Send("OK")
}

func AddToQueue(execute CommandFunc, args []string) {
	queued = append(queued, queuedCommand{
		execute: execute,
		args:    append([]string(nil), args...),
	})
}
