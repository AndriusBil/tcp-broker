package broker

import (
	"github.com/andriusbil/tcp-broker/logger"
	"net"
	"sync"
)

type Server struct {
	port        string
	listener    net.Listener
	Stream      chan string
	quit        chan bool
	log         logger.Logger
	Errors      chan error
	connHandler func(logger.Logger, chan string, net.Conn)
	mu          sync.Mutex
	Started     chan bool
}

func NewServer(
	port string,
	log logger.Logger,
	connHandler func(logger.Logger, chan string, net.Conn),
) *Server {
	return &Server{
		port:        port,
		Stream:      make(chan string),
		quit:        make(chan bool, 1),
		log:         log,
		Errors:      make(chan error),
		connHandler: connHandler,
		Started:     make(chan bool),
	}
}

func (s *Server) Start() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", s.port)
	if err != nil {
		s.Errors <- err
		return
	}

	l, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		s.Errors <- err
		return
	}

	s.mu.Lock()
	s.listener = l
	s.mu.Unlock()

	s.Started <- true

	defer l.Close()

	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			select {
			case <-s.quit:
				return
			default:
				s.log.Printf("%v", err)
			}
		}

		if err := conn.SetKeepAlive(true); err != nil {
			s.log.Printf("%v", err)
		}

		go s.connHandler(s.log, s.Stream, conn)
	}
}

func (s *Server) GetListener() net.Listener {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.listener
}

func (s *Server) Stop() {
	s.quit <- true
	if l := s.GetListener(); l != nil {
		l.Close()
	}
}
