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
	return &PublisherServer{
		port:   port,
		quit:   make(chan bool, 1),
		Stream: make(chan string),
	}
}

func handleInConnection(stream chan string, conn net.Conn) {
	content, err := io.ReadAll(conn)
	if err != nil {
		log.Print(err)
	}

	stream <- string(content)

	if err := conn.Close(); err != nil {
		log.Print(err)
	}
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
			}
		}

		if err := conn.SetKeepAlive(true); err != nil {
			log.Printf("%v", err)
		}

		go handleInConnection(ps.Stream, conn)
	}
}

func (ps *PublisherServer) Stop() error {
	ps.quit <- true
	ps.listener.Close()
	return nil
}
