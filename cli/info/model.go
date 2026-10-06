package info

import "fmt"

var Redis = RedisInfo{
	Port: 6379,
	Replication: Replication{
		Role: "master",
	},
}

type RedisInfo struct {
	Port        int
	Replication Replication
}

type Replication struct {
	Role string
}

func (r *Replication) String() string {
	return fmt.Sprintf("role:%s", r.Role)
}
