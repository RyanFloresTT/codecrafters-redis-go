package info

import "fmt"

var Redis = RedisInfo{
	Port: 6379,
	Replication: Replication{
		role: "master",
	},
}

type RedisInfo struct {
	Port        int
	Replication Replication
}

type Replication struct {
	role string
}

func (r *Replication) String() string {
	return fmt.Sprintf("role:%s", r.role)
}
