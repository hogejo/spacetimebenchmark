package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Connection struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
	buf    []byte
}

func newConnection(addr string) (*Connection, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("connect to server: %w", err)
	}
	tcpConn := conn.(*net.TCPConn)
	tcpConn.SetNoDelay(true)
	return &Connection{
		conn:   conn,
		reader: bufio.NewReaderSize(conn, 256),
		writer: bufio.NewWriterSize(conn, 256),
		buf:    make([]byte, 0, 128),
	}, nil
}

func (c *Connection) Close() error {
	return c.conn.Close()
}

func (c *Connection) sendTransfer(fromID, toID, amount uint64) (string, error) {
	c.buf = c.buf[:0]
	c.buf = append(c.buf, "transfer "...)
	c.buf = strconv.AppendUint(c.buf, fromID, 10)
	c.buf = append(c.buf, ' ')
	c.buf = strconv.AppendUint(c.buf, toID, 10)
	c.buf = append(c.buf, ' ')
	c.buf = strconv.AppendUint(c.buf, amount, 10)
	c.buf = append(c.buf, '\n')
	if _, err := c.writer.Write(c.buf); err != nil {
		return "", fmt.Errorf("write transfer: %w", err)
	}
	if err := c.writer.Flush(); err != nil {
		return "", fmt.Errorf("flush transfer: %w", err)
	}
	line, err := c.reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("read transfer response: %w", err)
	}
	return strings.TrimRight(line, "\r\n"), nil
}

func (c *Connection) sendGet(accountID uint64) (string, error) {
	c.buf = c.buf[:0]
	c.buf = append(c.buf, "get "...)
	c.buf = strconv.AppendUint(c.buf, accountID, 10)
	c.buf = append(c.buf, '\n')
	if _, err := c.writer.Write(c.buf); err != nil {
		return "", fmt.Errorf("write get: %w", err)
	}
	if err := c.writer.Flush(); err != nil {
		return "", fmt.Errorf("flush get: %w", err)
	}
	line, err := c.reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("read get response: %w", err)
	}
	return strings.TrimRight(line, "\r\n"), nil
}
