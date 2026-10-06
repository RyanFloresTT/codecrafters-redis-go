package replication

import (
	"encoding/base64"
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/cli/info"
	"github.com/codecrafters-io/redis-starter-go/resp"
)

func Psync(c resp.Connection, args []string) (resp.Value, error) {
	data, err := base64.StdEncoding.DecodeString(RDBContentsB64)
	if err != nil {
		fmt.Println("Error decoding RDB contents:", err)
	}

	return resp.CommandArray{
		resp.SimpleString(fmt.Sprintf("FULLRESYNC %s %d", info.Redis.Replication.MasterReplID, info.Redis.Replication.MasterReplOffset)),
		resp.File(string(data)),
	}, nil
}

const RDBContentsB64 = "UkVESVMwMDEx+glyZWRpcy12ZXIFNy4yLjD6CnJlZGlzLWJpdHPAQPoFY3RpbWXCbQi8ZfoIdXNlZC1tZW3CsMQQAPoIYW9mLWJhc2XAAP/wbjv+wP9aog=="
