package cli

import (
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/cli/info"
)

func ReplicaOf(args []string) error {
	masterParts := strings.Split(args[0], " ")

	masterHost := strings.TrimPrefix(masterParts[0], "\"")
	masterPort := strings.TrimSuffix(masterParts[1], "\"")

	info.Redis.Replication.MasterHost = masterHost
	info.Redis.Replication.MasterPort, _ = strconv.Atoi(masterPort)

	info.Redis.Replication.Role = "slave"

	return nil
}
