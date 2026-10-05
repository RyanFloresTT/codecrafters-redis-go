package helpers

import (
	"fmt"
	"net"
)

type Connection struct {
	net.Conn
}

type Value interface {
	SendTo(Connection) error
}

type SimpleString string
type BulkString string
type Integer int64
type Array []Value
type NullBulk struct{}
type NullArray struct{}
type Error string

func (message SimpleString) SendTo(c Connection) error {
	return c.Send(string(message))
}

func (message BulkString) SendTo(c Connection) error {
	return c.SendBulk(string(message))
}

func (number Integer) SendTo(c Connection) error {
	_, err := fmt.Fprintf(c.Conn, ":%d\r\n", number)
	return err
}

func (messages Array) SendTo(c Connection) error {
	if _, err := fmt.Fprintf(c.Conn, "*%d\r\n", len(messages)); err != nil {
		return err
	}
	for _, message := range messages {
		if err := message.SendTo(c); err != nil {
			return err
		}
	}
	return nil
}

func (NullBulk) SendTo(c Connection) error {
	return c.SendNull()
}

func (NullArray) SendTo(c Connection) error {
	return c.SendNullArray()
}

func (message Error) SendTo(c Connection) error {
	return c.SendError(string(message))
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
