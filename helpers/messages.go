package helpers

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
	_, err := fmt.Fprintf(c.Conn, "-%s\r\n", message)
	return err
}

func (c Connection) SendArray(messages []string) error {
	if _, err := fmt.Fprintf(c.Conn, "*%d\r\n", len(messages)); err != nil {
		return err
	}
	for _, message := range messages {
		if err := c.SendBulk(message); err != nil {
			return err
		}
	}
	return nil
}

func (c Connection) SendArrayLength(length int) error {
	_, err := fmt.Fprintf(c.Conn, ":%d\r\n", length)
	return err
}
