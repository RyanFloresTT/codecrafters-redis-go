package resp

import (
	"fmt"
	"net"
)

type Connection struct {
	net.Conn
}

func (c Connection) Send(message string) error {
	_, err := fmt.Fprintf(c.Conn, "+%s\r\n", message)
	return err
}

func (c Connection) SendBulk(message string) error {
	_, err := fmt.Fprintf(c.Conn, "$%d\r\n%s\r\n", len(message), message)
	return err
}

func (c Connection) SendNull() error {
	_, err := fmt.Fprint(c.Conn, "$-1\r\n")
	return err
}

func (c Connection) SendNullArray() error {
	_, err := fmt.Fprint(c.Conn, "*-1\r\n")
	return err
}

func (c Connection) SendError(message string) error {
	_, err := fmt.Fprintf(c.Conn, "-ERR %s\r\n", message)
	return err
}

func (c Connection) SendArray(messages []string) error {
	values := make(Array, len(messages))
	for index, message := range messages {
		values[index] = BulkString(message)
	}
	return values.SendTo(c)
}

func (c Connection) SendArrayLength(length int) error {
	return Integer(length).SendTo(c)
}

func (c Connection) SendFile(contents string) error {
	_, err := fmt.Fprintf(c.Conn, "$%d\r\n%s", len(contents), contents)
	return err
}
