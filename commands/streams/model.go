package streams

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

type entry struct {
	id     id
	values []kvp
}

func (item entry) RESP() resp.Value {
	fields := make(resp.Array, 0, len(item.values)*2)
	for _, pair := range item.values {
		fields = append(fields, resp.BulkString(pair.key), resp.BulkString(pair.value))
	}
	return resp.Array{
		resp.BulkString(item.id.String()),
		fields,
	}
}

type kvp struct {
	key   string
	value string
}

type id struct {
	ms  int64
	seq int
}

func (i id) String() string {
	return strconv.FormatInt(i.ms, 10) + "-" + strconv.Itoa(i.seq)
}

func GetLastID(key string) (id, bool) {
	stream := streamMap[key]
	if len(stream) == 0 {
		return id{0, 0}, false
	}
	lastItem := stream[len(stream)-1]
	return lastItem.id, true
}
