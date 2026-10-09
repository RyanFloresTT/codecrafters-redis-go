package numbers

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/commands/strings"
	"github.com/codecrafters-io/redis-starter-go/resp"
)

func Incr(c resp.Connection, args []string) resp.Value {
	if len(args) != 1 {
		return resp.Error("wrong number of arguments for 'incr' command")
	}
	entry, ok := strings.GetEntry(args[0])

	if !ok {
		entry := strings.Entry()
		entry.Value = "1"
		strings.GetMap()[args[0]] = entry
		return resp.Integer(1)
	}

	value, err := strconv.Atoi(entry.Value)

	if err != nil {
		return resp.Error("value is not an integer or out of range")
	}

	data := resp.Integer(value + 1)
	entry.Value = strconv.Itoa(value + 1)

	strings.GetMap()[args[0]] = entry

	return data
}
