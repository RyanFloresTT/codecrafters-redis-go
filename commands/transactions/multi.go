package transactions

import (
	"net"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

type clientTransaction struct {
	commands []queuedCommand
}

var (
	clientsMu sync.Mutex
	clients   = make(map[net.Addr]*clientTransaction)
)

type CommandFunc func(helpers.Connection, []string) (helpers.Value, error)

type queuedCommand struct {
	execute CommandFunc
	args    []string
}

func Multi(c helpers.Connection, args []string) (helpers.Value, error) {
	clientsMu.Lock()
	clients[c.RemoteAddr()] = &clientTransaction{}
	clientsMu.Unlock()
	return helpers.SimpleString("OK"), nil
}

func Discard(c helpers.Connection, args []string) (helpers.Value, error) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	if _, ok := clients[c.RemoteAddr()]; ok {
		delete(clients, c.RemoteAddr())
		return helpers.SimpleString("OK"), nil
	}

	return helpers.Error("DISCARD without MULTI"), nil
}

func IsActive(c helpers.Connection) bool {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	_, ok := clients[c.RemoteAddr()]
	return ok
}

func AddToQueue(c helpers.Connection, execute CommandFunc, args []string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	transaction := clients[c.RemoteAddr()]
	if transaction == nil {
		return
	}
	transaction.commands = append(transaction.commands, queuedCommand{
		execute: execute,
		args:    append([]string(nil), args...),
	})
}
