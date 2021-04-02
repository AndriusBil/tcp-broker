package broker

import "github.com/andriusbil/tcp-broker/logger"

type BrokerServer struct {
	publishersPort string
	consumersPort  string

	consumerServer  *ConsumerServer
	publisherServer *PublisherServer

	quit chan bool
}

func NewBrokerServer(publishersPort string, consumersPort string, log logger.Logger) *BrokerServer {
	return &BrokerServer{
		publishersPort:  publishersPort,
		consumersPort:   consumersPort,
		consumerServer:  NewConsumerServer(consumersPort, log),
		publisherServer: NewPublisherServer(publishersPort, log),
		quit:            make(chan bool, 1),
	}
}

func (bs *BrokerServer) Start() {

	go bs.consumerServer.Start()
	go bs.publisherServer.Start()

	for {
		select {
		case msg := <-bs.publisherServer.Stream:
			bs.consumerServer.Stream <- msg
		case <-bs.quit:
		default:
		}
	}
}

func (bs *BrokerServer) Stop() {
	bs.quit <- true
	bs.consumerServer.Stop()
	bs.publisherServer.Stop()
}
