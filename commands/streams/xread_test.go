package streams

import (
	"bytes"
	"net"
	"testing"

	"github.com/codecrafters-io/redis-starter-go/helpers"
)

type xreadTestConnection struct {
	net.Conn
	output bytes.Buffer
}

func (connection *xreadTestConnection) Write(data []byte) (int, error) {
	return connection.output.Write(data)
}

func TestXRead(t *testing.T) {
	original := streamMap
	t.Cleanup(func() { streamMap = original })
	first := entry{id: id{ms: 0, seq: 1}, values: []kvp{{key: "temperature", value: "88"}}}
	later := entry{id: id{ms: 0, seq: 3}, values: first.values}
	nextTimestamp := entry{id: id{ms: 1, seq: 0}, values: first.values}
	const prefix = "*1\r\n*2\r\n$9\r\npineapple\r\n"
	const firstRESP = "*2\r\n$3\r\n0-1\r\n*2\r\n$11\r\ntemperature\r\n$2\r\n88\r\n"
	const laterRESP = "*2\r\n$3\r\n0-3\r\n*2\r\n$11\r\ntemperature\r\n$2\r\n88\r\n"
	const nextRESP = "*2\r\n$3\r\n1-0\r\n*2\r\n$11\r\ntemperature\r\n$2\r\n88\r\n"
	tests := []struct {
		name    string
		entries []entry
		args    []string
		want    string
	}{
		{"reported nesting", []entry{first}, []string{"streams", "pineapple", "0-0"}, prefix + "*1\r\n" + firstRESP},
		{"exclusive bound", []entry{first, later, nextTimestamp}, []string{"streams", "pineapple", "0-1"}, prefix + "*2\r\n" + laterRESP + nextRESP},
		{"nonexistent bound", []entry{first, later}, []string{"streams", "pineapple", "0-2"}, prefix + "*1\r\n" + laterRESP},
		{"numeric comparison", []entry{first, later, nextTimestamp}, []string{"streams", "pineapple", "0-10"}, prefix + "*1\r\n" + nextRESP},
		{"no newer entries", []entry{first}, []string{"streams", "pineapple", "0-1"}, "*-1\r\n"},
		{"missing stream", nil, []string{"streams", "pineapple", "0-0"}, "*-1\r\n"},
		{"missing arguments", nil, nil, "-ERR wrong number of arguments for 'xread' command\r\n"},
		{"unpaired arguments", nil, []string{"streams", "pineapple", "mango", "0-0"}, "-ERR wrong number of arguments for 'xread' command\r\n"},
		{"invalid keyword", nil, []string{"invalid", "pineapple", "0-0"}, "-ERR syntax error\r\n"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			streamMap = map[string][]entry{"pineapple": test.entries}
			connection := &xreadTestConnection{}
			response, err := XRead(helpers.Connection{Conn: connection}, test.args)
			if err != nil {
				t.Fatal(err)
			}
			if err := response.SendTo(helpers.Connection{Conn: connection}); err != nil {
				t.Fatal(err)
			}
			if got := connection.output.String(); got != test.want {
				t.Fatalf("got %q, want %q", got, test.want)
			}
		})
	}
}

func TestXReadMultipleStreams(t *testing.T) {
	original := streamMap
	t.Cleanup(func() { streamMap = original })
	streamMap = map[string][]entry{
		"pineapple": {
			{id: id{ms: 0, seq: 1}, values: []kvp{{key: "temperature", value: "88"}}},
			{id: id{ms: 0, seq: 3}, values: []kvp{{key: "temperature", value: "90"}}},
		},
		"0-9": {
			{id: id{ms: 0, seq: 1}, values: []kvp{{key: "temperature", value: "77"}}},
		},
	}
	const pineapple = "*2\r\n$9\r\npineapple\r\n*1\r\n*2\r\n$3\r\n0-3\r\n*2\r\n$11\r\ntemperature\r\n$2\r\n90\r\n"
	const numericKey = "*2\r\n$3\r\n0-9\r\n*1\r\n*2\r\n$3\r\n0-1\r\n*2\r\n$11\r\ntemperature\r\n$2\r\n77\r\n"
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"paired IDs and numeric key", []string{"streams", "pineapple", "0-9", "0-1", "0-0"}, "*2\r\n" + pineapple + numericKey},
		{"request order", []string{"streams", "0-9", "pineapple", "0-0", "0-1"}, "*2\r\n" + numericKey + pineapple},
		{"empty first stream", []string{"streams", "missing", "pineapple", "0-0", "0-1"}, "*1\r\n" + pineapple},
		{"empty last stream", []string{"streams", "pineapple", "0-9", "0-1", "0-1"}, "*1\r\n" + pineapple},
		{"all streams empty", []string{"streams", "pineapple", "0-9", "0-3", "0-1"}, "*-1\r\n"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			connection := &xreadTestConnection{}
			response, err := XRead(helpers.Connection{Conn: connection}, test.args)
			if err != nil {
				t.Fatal(err)
			}
			if err := response.SendTo(helpers.Connection{Conn: connection}); err != nil {
				t.Fatal(err)
			}
			if got := connection.output.String(); got != test.want {
				t.Fatalf("got %q, want %q", got, test.want)
			}
		})
	}
}
