package numbers

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/commands/strings"
	"github.com/codecrafters-io/redis-starter-go/helpers"
)

func Incr(c helpers.Connection, args []string) (helpers.Value, error) {
	if len(args) != 1 {
		return helpers.Error("wrong number of arguments for 'incr' command"), nil
	}
	entry, ok := strings.GetEntry(args[0])

	if !ok {
		entry := strings.Entry()
		entry.Value = "1"
		strings.GetMap()[args[0]] = entry
		return helpers.Integer(1), nil
	}

	value, err := strconv.Atoi(entry.Value)

	if err != nil {
		return helpers.Error("value is not an integer or out of range"), nil
	}

	data := helpers.Integer(value + 1)
	entry.Value = strconv.Itoa(value + 1)

	strings.GetMap()[args[0]] = entry

	return data, nil
}
