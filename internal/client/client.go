package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"net"

	"word_of_wisdom/internal/pkg/data"
	"word_of_wisdom/internal/pkg/pow"
)

const (
	emptyString = ""
)

type Client struct {
	conn net.Conn
}

func NewClient(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("client.Start: %w", err)
	}

	return &Client{conn: conn}, nil
}

func (c *Client) Start() error {
	defer c.conn.Close()

	for {
		res, err := c.handleConnection()
		if err != nil {
			return fmt.Errorf("client.Start: %w", err)
		}
		fmt.Println("Word of Wisdom:", res)
	}
}

func (c *Client) handleConnection() (string, error) {
	err := c.sendData(data.Data{
		Key:   data.Challenge,
		Value: emptyString,
	})
	if err != nil {
		return emptyString, fmt.Errorf("client.handleConnection: %w", err)
	}

	reader := bufio.NewReader(c.conn)
	dataStr, err := readData(reader)
	if err != nil {
		return emptyString, fmt.Errorf("client.handleConnection - read data: %w", err)
	}

	hashcash, err := hashing(dataStr)
	if err != nil {
		return emptyString, fmt.Errorf("client.handleConnection - hashing: %w", err)
	}

	byteData, err := json.Marshal(hashcash)
	if err != nil {
		return emptyString, fmt.Errorf("client.handleConnection - marshal hashcash: %w", err)
	}

	err = c.sendData(data.Data{
		Key:   data.Response,
		Value: string(byteData),
	})
	if err != nil {
		return emptyString, fmt.Errorf("client.handleConnection - send request: %w", err)
	}

	dataStr, err = readData(reader)
	if err != nil {
		return emptyString, fmt.Errorf("client.handleConnection - read data: %w", err)
	}

	d, err := data.Parse(dataStr)
	if err != nil {
		return emptyString, fmt.Errorf("client.handleConnection - parse data: %w", err)
	}
	return d.Value, nil
}

func hashing(dataStr string) (pow.HashcashData, error) {
	var hashcash pow.HashcashData
	d, err := data.Parse(dataStr)
	if err != nil {
		return hashcash, fmt.Errorf("client.handleConnection - parse data: %w", err)
	}
	err = json.Unmarshal([]byte(d.Value), &hashcash)
	if err != nil {
		return hashcash, fmt.Errorf("client.handleConnection - parse hashcash: %w", err)
	}

	hashcash, err = hashcash.ComputeHashcash(math.MaxInt64)
	if err != nil {
		return hashcash, fmt.Errorf("client.handleConnection - compute hashcash: %w", err)
	}

	return hashcash, nil
}

func readData(reader *bufio.Reader) (string, error) {
	return reader.ReadString('\n')
}

func (c *Client) sendData(msg data.Data) error {
	dataStr := fmt.Sprintf("%s\n", msg.String())
	_, err := c.conn.Write([]byte(dataStr))
	return err
}
