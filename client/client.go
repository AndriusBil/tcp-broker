package client

import (
	"bufio"
	"log"
	"net"
	"strings"
)

type Client struct {
	url  string
	port string
	conn net.Conn
}

func New(url string, port string) *Client {
	return &Client{url: url, port: port}
}

func (c *Client) SendMessage(message string) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", c.url+c.port)
	if err != nil {
		return err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}

	if err := conn.SetKeepAlive(true); err != nil {
		return err
	}

	defer conn.Close()

	if _, err := conn.Write([]byte(message)); err != nil {
		return err
	}

	return nil
}

func listenMessages(conn net.Conn, fn func(string)) {
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')

		if err != nil {
			log.Printf("%v", err)
			return
		}

		go fn(strings.Trim(strings.TrimSpace(msg), "\n"))
	}
}

func (c *Client) Subscribe(fn func(string)) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", c.url+c.port)
	if err != nil {
		log.Printf("%v", err)
		return err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Printf("%v", err)
		return err
	}

	conn.SetKeepAlive(true)

	c.conn = conn

	go listenMessages(conn, fn)

	return nil
}

func (c *Client) Stop() {
	if c.conn != nil {
		c.conn.Close()
	}
}
