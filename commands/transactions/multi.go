package transactions

import (
	"net"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/commands/optimistic_locking/state"
	"github.com/codecrafters-io/redis-starter-go/resp"
)

type clientTransaction struct {
	commands []queuedCommand
}

var (
	clientsMu sync.Mutex
	clients   = make(map[net.Addr]*clientTransaction)
)

type CommandFunc func(resp.Connection, []string) resp.Value

type queuedCommand struct {
	execute CommandFunc
	args    []string
}

func Multi(c resp.Connection, args []string) resp.Value {
	clientsMu.Lock()
	clients[c.RemoteAddr()] = &clientTransaction{}
	clientsMu.Unlock()
	return resp.SimpleString("OK")
}

func Discard(c resp.Connection, args []string) resp.Value {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	state.WatchedKeys.Clear()

	if _, ok := clients[c.RemoteAddr()]; ok {
		delete(clients, c.RemoteAddr())
		return resp.SimpleString("OK")
	}

	return resp.Error("DISCARD without MULTI")
}

func IsActive(c resp.Connection) bool {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	_, ok := clients[c.RemoteAddr()]
	return ok
}

func AddToQueue(c resp.Connection, execute CommandFunc, args []string) {
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
