package broker

type BrokerServer struct {
	publishersPort string
	consumersPort  string

	consumerServer  *ConsumerServer
	publisherServer *PublisherServer

	quit chan bool
}

func NewBrokerServer(publishersPort string, consumersPort string) *BrokerServer {
	return &BrokerServer{
		publishersPort: publishersPort,
		consumersPort:  consumersPort,
		quit:           make(chan bool, 1),
	}
}

func (bs *BrokerServer) Start() {
	bs.consumerServer = NewConsumerServer(bs.consumersPort)
	bs.publisherServer = NewPublisherServer(bs.publishersPort)

	go bs.consumerServer.Start()
	go bs.publisherServer.Start()

	for {
		select {
		case msg := <-bs.publisherServer.Stream:
			bs.consumerServer.Stream <- msg
		case <-bs.quit:
		default:
			continue
		}
	}
}

func (bs *BrokerServer) Stop() {
	bs.quit <- true
	bs.consumerServer.Stop()
	bs.publisherServer.Stop()
}
