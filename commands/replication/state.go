package replication

import "github.com/codecrafters-io/redis-starter-go/resp"

var replicas []resp.Connection

func GetReplicas() []resp.Connection {
	return replicas
}

func AddReplica(c resp.Connection) {
	replicas = append(replicas, c)
}
