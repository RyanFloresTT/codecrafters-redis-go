package numbers

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/commands/strings"
	"github.com/codecrafters-io/redis-starter-go/helpers"
)

func Incr(c helpers.Connection, args []string) error {
	// Implementation for INCR command goes here
	dictionary := strings.GetMap()
	entry := dictionary[args[0]]

	value, err := strconv.Atoi(entry.Value)
	if err != nil {
		return c.SendError("ERR value is not an integer or out of range")
	}

	data := helpers.Integer(value + 1)

	entry.Value = strconv.Itoa(value + 1)
	dictionary[args[0]] = entry
	return data.SendTo(c)
}
