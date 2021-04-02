package main

import (
	"github.com/andriusbil/tcp-broker/broker"
	"log"
	"os"
)

func main() {
	server := broker.NewBrokerServer(
		os.Getenv("PUBLISHERS_PORT"),
		os.Getenv("CONSUMERS_PORT"),
		log.Default(),
	)
	server.Start()
}
