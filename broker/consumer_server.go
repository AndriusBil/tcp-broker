package broker

import (
	"log"
	"net"
)

type ConsumerServer struct {
	port     string
	listener net.Listener
	Stream   chan string
	quit     chan bool
}

func NewConsumerServer(port string) *ConsumerServer {
	return &ConsumerServer{
		port:   port,
		Stream: make(chan string),
		quit:   make(chan bool, 1),
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

func (cs *ConsumerServer) Start() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", cs.port)
	if err != nil {
		return err
	}

	l, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	cs.listener = l

	defer l.Close()

	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			select {
			case <-cs.quit:
				return err
			default:
				log.Printf("%v", err)
			}
		}

		if err := conn.SetKeepAlive(true); err != nil {
			log.Printf("%v", err)
		}

		go handleOutConnection(cs.Stream, conn)
	}
}

func (cs *ConsumerServer) Stop() error {
	cs.quit <- true
	cs.listener.Close()
	return nil
}
