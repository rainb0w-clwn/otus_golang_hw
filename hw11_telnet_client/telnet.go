package main

import (
	"errors"
	"io"
	"net"
	"time"
)

var (
	ErrConnectionEmpty          = errors.New("telnet connection not opened")
	ErrConnectionClosedByClient = errors.New("...Connection was closed by client")
	ErrConnectionClosedByPeer   = errors.New("...Connection was closed by peer")
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}
type telnetClient struct {
	address string
	timeout time.Duration
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
}

func (c *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}
	errLog.Println("...Connected to " + c.address)
	c.conn = conn
	return err
}

func (c *telnetClient) Send() error {
	if c.conn == nil {
		return ErrConnectionEmpty
	}
	if _, err := io.Copy(c.conn, c.in); err != nil {
		return err
	}
	errLog.Println(ErrConnectionClosedByClient)
	return nil
}

func (c *telnetClient) Receive() error {
	if c.conn == nil {
		return ErrConnectionEmpty
	}
	if _, err := io.Copy(c.out, c.conn); err != nil {
		return err
	}
	errLog.Println(ErrConnectionClosedByPeer)
	return nil
}

func (c *telnetClient) Close() error {
	if c.conn == nil {
		return ErrConnectionEmpty
	}
	err := c.conn.Close()
	if err == nil {
		c.conn = nil
	}
	return err
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}
