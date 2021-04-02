package broker

import (
	"github.com/andriusbil/tcp-broker/logger"
	"io"
	"net"
)

type PublisherServer struct {
	port     string
	listener net.Listener
	quit     chan bool
	Stream   chan string
	log      logger.Logger
	Errors   chan error
}

func NewPublisherServer(port string, log logger.Logger) *PublisherServer {
	return &PublisherServer{
		port:   port,
		quit:   make(chan bool, 1),
		Stream: make(chan string),
		log:    log,
		Errors: make(chan error),
	}
}

func handleInConnection(log logger.Logger, stream chan string, conn net.Conn) {
	content, err := io.ReadAll(conn)
	if err != nil {
		log.Print(err)
	}

	stream <- string(content)

	if err := conn.Close(); err != nil {
		log.Print(err)
	}
}

func (ps *PublisherServer) Start() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ps.port)
	if err != nil {
		ps.Errors <- err
		return
	}

	l, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		ps.Errors <- err
		return
	}

	ps.listener = l

	defer l.Close()

	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			select {
			case <-ps.quit:
				ps.Errors <- err
				return
			default:
				ps.log.Printf("%v", err)
			}
		}

		if err := conn.SetKeepAlive(true); err != nil {
			ps.log.Printf("%v", err)
		}

		go handleInConnection(ps.log, ps.Stream, conn)
	}
}

func (ps *PublisherServer) Stop() error {
	ps.quit <- true
	if ps.listener != nil {
		ps.listener.Close()
	}
	return nil
}
