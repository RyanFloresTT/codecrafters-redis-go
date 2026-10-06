package replication

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/cli/info"
	"github.com/codecrafters-io/redis-starter-go/resp"
)

func Psync(c resp.Connection, args []string) (resp.Value, error) {
	return resp.SimpleString(fmt.Sprintf("FULLRESYNC %s %d", info.Redis.Replication.MasterReplID, info.Redis.Replication.MasterReplOffset)), nil
}
