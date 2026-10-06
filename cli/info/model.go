package info

import "fmt"

var Redis = RedisInfo{
	Port: 6379,
	Replication: Replication{
		Role:             "master",
		MasterReplID:     "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
		MasterReplOffset: "0",
	},
}

type RedisInfo struct {
	Port        int
	Replication Replication
}

type Replication struct {
	Role             string
	MasterReplID     string
	MasterReplOffset string
}

func (r *Replication) String() string {
	return fmt.Sprintf("role:%s\nmaster_replid:%s\nmaster_repl_offset:%s\n", r.Role, r.MasterReplID, r.MasterReplOffset)
}
