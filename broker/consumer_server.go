package broker

import (
	"github.com/andriusbil/tcp-broker/logger"
	"net"
)

func handleOutConnection(log logger.Logger, in chan string, conn net.Conn) {
	for msg := range in {
		if _, err := conn.Write([]byte(msg + string('\n'))); err != nil {
			conn.Close()
			log.Printf("%v", err)
			return
		}
	}
}

func NewConsumerServer(port string, log logger.Logger) *Server {
	return NewServer(port, log, handleOutConnection)
}
