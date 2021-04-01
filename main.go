package main

import (
	"github.com/andriusbil/tcp-broker/broker"
	"os"
)

func main() {
	server := broker.NewBrokerServer(os.Getenv("PUBLISHERS_PORT"), os.Getenv("CONSUMERS_PORT"))
	defer server.Stop()
	server.Start()
}
