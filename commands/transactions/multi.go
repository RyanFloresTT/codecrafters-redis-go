package transactions

import (
	"github.com/codecrafters-io/redis-starter-go/helpers"
)

var IsTransaction bool

type CommandFunc func(helpers.Connection, []string) error

type queuedCommand struct {
	execute CommandFunc
	args    []string
}

var queued []queuedCommand

func Multi(c helpers.Connection, args []string) error {
	IsTransaction = true
	return c.Send("OK")
}

func AddToQueue(execute CommandFunc, args []string) {
	queued = append(queued, queuedCommand{
		execute: execute,
		args:    append([]string(nil), args...),
	})
}
