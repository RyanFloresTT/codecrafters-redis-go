package commands

import "github.com/codecrafters-io/redis-starter-go/resp"

type Command struct {
	Name              string
	Execute           func(connection resp.Connection, args []string) (resp.Value, error)
	Args              []string
	IsExemptFromQueue bool
	Replicates        bool
}

type CLICommand struct {
	Name              string
	Execute           func(args []string) error
	Args              []string
	IsExemptFromQueue bool
}

type ReplicaCommand struct {
	CommandRESP resp.Array
}

func (r ReplicaCommand) SendTo(connection resp.Connection) error {
	return r.CommandRESP.SendTo(connection)
}

func (c Command) ToReplicaCommand() ReplicaCommand {
	values := make(resp.Array, len(c.Args)+1)
	values[0] = resp.BulkString(c.Name)

	for index, arg := range c.Args {
		values[index+1] = resp.BulkString(arg)
	}

	return ReplicaCommand{values}
}
