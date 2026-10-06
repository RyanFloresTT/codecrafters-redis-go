package resp

import "fmt"

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
