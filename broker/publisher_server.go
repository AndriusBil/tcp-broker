package broker

import (
	"github.com/andriusbil/tcp-broker/logger"
	"io"
	"net"
)

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

func NewPublisherServer(port string, log logger.Logger) *Server {
	return NewServer(port, log, handleInConnection)
}
