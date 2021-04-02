package broker

import (
	"github.com/andriusbil/tcp-broker/logger"
	"net"
)

type ConsumerServer struct {
	port     string
	listener net.Listener
	Stream   chan string
	quit     chan bool
	log      logger.Logger
	Errors   chan error
}

func NewConsumerServer(port string, log logger.Logger) *ConsumerServer {
	return &ConsumerServer{
		port:   port,
		Stream: make(chan string),
		quit:   make(chan bool, 1),
		log:    log,
		Errors: make(chan error),
	}
}

func handleOutConnection(in chan string, conn net.Conn) {
	for {
		msg := <-in
		if _, err := conn.Write([]byte(msg + string('\n'))); err != nil {
			conn.Close()
		}
	}
}

func (cs *ConsumerServer) Start() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", cs.port)
	if err != nil {
		cs.Errors <- err
		return
	}

	l, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		cs.Errors <- err
		return
	}

	cs.listener = l

	defer l.Close()

	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			select {
			case <-cs.quit:
				cs.Errors <- err
				return
			default:
				cs.log.Printf("%v", err)
			}
		}

		if err := conn.SetKeepAlive(true); err != nil {
			cs.log.Printf("%v", err)
		}

		go handleOutConnection(cs.Stream, conn)
	}
}

func (cs *ConsumerServer) Stop() error {
	cs.quit <- true
	if cs.listener != nil {
		cs.listener.Close()
	}
	return nil
}
