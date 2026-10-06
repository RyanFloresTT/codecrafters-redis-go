package cli

import "github.com/codecrafters-io/redis-starter-go/cli/info"

func ReplicaOf(args []string) error {
	// masterHost := args[0]
	// masterPort := args[1]

	// info.Redis.ReplicaOfHost = masterHost
	// info.Redis.ReplicaOfPort = masterPort

	info.Redis.Replication.Role = "slave"

	return nil
}
