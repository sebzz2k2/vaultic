package server

import (
	"bufio"
	"errors"
	"net"
	"time"

	"github.com/sebzz2k2/vaultic/internal/protocol/lexer"
	"github.com/sebzz2k2/vaultic/internal/storage"
)

type Client struct {
	conn   net.Conn
	engine storage.StorageEngine
	config *Config
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewClient(conn net.Conn, config *Config, engine storage.StorageEngine) *Client {
	return &Client{
		conn:   conn,
		engine: engine,
		config: config,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
	}
}

func (c *Client) read() ([]byte, error) {
	b := make([]byte, 1024)
	bn, err := c.reader.Read(b)
	if err != nil {
		return nil, errors.New("failed to read from client: " + err.Error())
	}
	return b[:bn], nil
}

func (c *Client) Handle() error {
	defer c.writer.Flush()

	for {
		if err := c.writeMessage("> "); err != nil {
			return err
		}
		buff, err := c.read()
		if err != nil {
			return err
		}
		tkns := lexer.Tokenize(string(buff))
		val, err := c.engine.Protocol.ProcessCommand(tkns)
		if err != nil {
			if err := c.writeMessage(err.Error() + "\n"); err != nil {
				return err
			}
		} else {
			if err := c.writeMessage(val + "\n"); err != nil {
				return err
			}
		}
		if c.config.ReadTimeout > 0 {
			c.conn.SetReadDeadline(time.Now().Add(c.config.ReadTimeout))
		}
		if c.config.WriteTimeout > 0 {
			c.conn.SetWriteDeadline(time.Now().Add(c.config.WriteTimeout))
		}
	}
}

func (c *Client) writeMessage(message string) error {
	_, err := c.writer.WriteString(message)
	if err != nil {
		return err
	}
	return c.writer.Flush()
}
