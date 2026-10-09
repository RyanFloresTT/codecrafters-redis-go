package info

import (
	"fmt"
)

var Redis = RedisInfo{
	Port: 6379,
	Replication: Replication{
		Role:             "master",
		MasterReplID:     "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
		MasterReplOffset: 0,
		MasterHost:       "127.0.0.1",
		MasterPort:       6379,
	},
}

type RedisInfo struct {
	Port        int
	Replication Replication
}

type Replication struct {
	Role             string
	MasterReplID     string
	MasterReplOffset int
	MasterHost       string
	MasterPort       int
}

func (r *Replication) String() string {
	return fmt.Sprintf("role:%s\nmaster_replid:%s\nmaster_repl_offset:%d\n", r.Role, r.MasterReplID, r.MasterReplOffset)
}

func (r *RedisInfo) AddToOffset(n int) {
	r.Replication.MasterReplOffset += n
}

func (r *RedisInfo) GetOffset() int {
	return r.Replication.MasterReplOffset
}
