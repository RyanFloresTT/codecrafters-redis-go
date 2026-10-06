package cli

import (
	"fmt"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/cli/info"
	"github.com/codecrafters-io/redis-starter-go/resp"
)

func Port(args []string) (resp.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("--port requires one value")
	}
	portNum, err := strconv.Atoi(args[0])
	if err != nil || portNum <= 0 || portNum > 65535 {
		return nil, fmt.Errorf("invalid port: %s", args[0])
	}

	info.Redis.Port = portNum
	return nil, nil
}
