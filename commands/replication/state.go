package replication

import (
	"sync"
	"time"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

var replicaStateMu sync.Mutex
var replicaStateChanged = sync.NewCond(&replicaStateMu)

var replicas []resp.Connection

func GetReplicas() []resp.Connection {
	replicaStateMu.Lock()
	defer replicaStateMu.Unlock()

	return append([]resp.Connection(nil), replicas...)
}

func AddReplica(c resp.Connection) {
	replicaStateMu.Lock()
	defer replicaStateMu.Unlock()

	replicas = append(replicas, c)
	replicaOffsets[c] = 0
	replicaStateChanged.Broadcast()
}

var replicaOffsets = make(map[resp.Connection]int)

func GetReplicaOffset(c resp.Connection) int {
	replicaStateMu.Lock()
	defer replicaStateMu.Unlock()

	return replicaOffsets[c]
}

func SetReplicaOffset(c resp.Connection, offset int) {
	replicaStateMu.Lock()
	defer replicaStateMu.Unlock()

	replicaOffsets[c] = offset
	replicaStateChanged.Broadcast()
}

func waitForReplicaAcks(targetOffset, required int, deadline time.Time) int {
	remaining := time.Until(deadline)
	if remaining < 0 {
		remaining = 0
	}
	timer := time.AfterFunc(remaining, func() {
		replicaStateMu.Lock()
		replicaStateChanged.Broadcast()
		replicaStateMu.Unlock()
	})
	defer timer.Stop()

	replicaStateMu.Lock()
	defer replicaStateMu.Unlock()

	if required <= 0 {
		return 0
	}

	for {
		caughtUp := 0
		for _, c := range replicas {
			if replicaOffsets[c] >= targetOffset {
				caughtUp++
			}
		}

		if caughtUp >= required {
			return caughtUp
		}

		if !time.Now().Before(deadline) {
			return caughtUp
		}
		replicaStateChanged.Wait()
	}
}
