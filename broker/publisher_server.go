package broker

import (
	"io"
	"log"
	"net"
)

type PublisherServer struct {
	port     string
	listener net.Listener
	quit     chan bool
	Stream   chan string
}

func NewPublisherServer(port string) *PublisherServer {
	server := PublisherServer{
		port: port,
	}
	server.quit = make(chan bool, 1)
	server.Stream = make(chan string)
	return &server
}

func (ps *PublisherServer) Start() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ps.port)
	if err != nil {
		return err
	}

	l, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	ps.listener = l

	defer l.Close()

	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			select {
			case <-ps.quit:
				return nil
			default:
				log.Printf("%v", err)
				continue
			}
		}

		if err := conn.SetKeepAlive(true); err != nil {
			log.Printf("%v", err)
		}

		go func(c net.Conn) {
			content, err := io.ReadAll(c)
			if err != nil {
				log.Print(err)
			}

			ps.Stream <- string(content)

			if err := c.Close(); err != nil {
				log.Print(err)
			}
		}(conn)
	}
}

func (ps *PublisherServer) Stop() error {
	ps.quit <- true
	ps.listener.Close()
	return nil
}
