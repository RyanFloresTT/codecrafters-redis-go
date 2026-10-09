package replication

import (
	"strconv"
	"time"

	"github.com/codecrafters-io/redis-starter-go/cli/info"
	"github.com/codecrafters-io/redis-starter-go/resp"
)

func Wait(c resp.Connection, args []string) (resp.Value, error) {
	numReplicas, _ := strconv.Atoi(args[0])
	timeout, _ := strconv.Atoi(args[1])

	targetOffset := info.Redis.GetOffset()
	deadline := time.Now().Add(time.Duration(timeout) * time.Millisecond)
	getReplicaOffsets()

	return resp.Integer(waitForReplicaAcks(targetOffset, numReplicas, deadline)), nil
}

func getReplicaOffsets() {
	for _, c := range GetReplicas() {
		sendGetAck(c)
	}
}

func sendGetAck(c resp.Connection) {
	response := resp.Array([]resp.Value{
		resp.BulkString("REPLCONF"),
		resp.BulkString("GETACK"),
		resp.BulkString("*"),
	})
	if err := response.SendTo(c); err != nil {
		return
	}
}
