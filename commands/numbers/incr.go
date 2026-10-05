package numbers

import (
	"fmt"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/commands/strings"
	"github.com/codecrafters-io/redis-starter-go/helpers"
)

func Incr(c helpers.Connection, args []string) error {
	entry, ok := strings.GetEntry(args[0])

	if !ok {
		entry := strings.Entry()
		entry.Value = "1"
		strings.GetMap()[args[0]] = entry
		return helpers.Integer(1).SendTo(c)
	}

	value, err := strconv.Atoi(entry.Value)

	if err != nil {
		return c.SendError("value is not an integer or out of range")
	}

	fmt.Println(value)

	data := helpers.Integer(value + 1)
	entry.Value = strconv.Itoa(value + 1)

	strings.GetMap()[args[0]] = entry

	return data.SendTo(c)
}
