package test

import (
	"github.com/andriusbil/tcp-broker/broker"
	"github.com/andriusbil/tcp-broker/client"
	"github.com/andriusbil/tcp-broker/logger"
	"github.com/stretchr/testify/assert"
	"log"
	"strings"
	"testing"
	"time"
)

func TestBroker(t *testing.T) {
	t.Run("Broker server", func(t *testing.T) {
		t.Run("it should notify consumers with new messages from publishers", func(t *testing.T) {
			message := "Test123"

			incomingMessages := logger.SliceWriter{}

			pp, cp := ":3333", ":3334"
			server := broker.NewBrokerServer(pp, cp, log.Default())
			go server.Start()
			defer server.Stop()

			publisher := client.New("localhost", pp)
			publisher2 := client.New("localhost", pp)
			consumer := client.New("localhost", cp)
			consumer2 := client.New("localhost", cp)

			<-server.Started

			err := consumer.Subscribe(func(msg string) {
				incomingMessages.Write([]byte(msg))
			})

			assert.Nil(t, err)
			defer consumer.Stop()

			err = consumer2.Subscribe(func(msg string) {
				incomingMessages.Write([]byte(msg))
			})
			assert.Nil(t, err)
			defer consumer2.Stop()

			publisher.SendMessage(message)
			publisher.SendMessage(message)
			publisher2.SendMessage(message)
			publisher2.SendMessage(message)

			assert.Eventually(t, func() bool {
				return assert.Contains(t, incomingMessages.Get(), message)
			}, time.Second, 10*time.Millisecond)

			assert.Eventually(t, func() bool {
				return assert.Len(t, incomingMessages.Get(), 4)
			}, time.Second, 10*time.Millisecond)
		})

		t.Run("it should notify new consumers with previously published messages", func(t *testing.T) {
			message := "Test123"

			incomingMessages := logger.SliceWriter{}

			pp, cp := ":3335", ":3336"
			server := broker.NewBrokerServer(pp, cp, log.Default())
			go server.Start()
			defer server.Stop()

			publisher := client.New("localhost", pp)
			consumer := client.New("localhost", cp)

			publisher.SendMessage(message)
			publisher.SendMessage(message)
			publisher.SendMessage(message)
			publisher.SendMessage(message)

			<-server.Started

			// Subscribe after push
			err := consumer.Subscribe(func(msg string) {
				incomingMessages.Write([]byte(msg))
			})
			assert.Nil(t, err)
			defer consumer.Stop()

			assert.Eventually(t, func() bool {
				return assert.Contains(t, incomingMessages.Get(), message)
			}, time.Second, 10*time.Millisecond)

			assert.Eventually(t, func() bool {
				return assert.Len(t, incomingMessages.Get(), 4)
			}, time.Second, 10*time.Millisecond)
		})

		t.Run("it should continue serving consumers if some of them stop", func(t *testing.T) {
			message := "Test123"

			incomingMessages := logger.SliceWriter{}

			pp, cp := ":3337", ":3338"
			server := broker.NewBrokerServer(pp, cp, log.Default())
			go server.Start()
			defer server.Stop()

			publisher := client.New("localhost", pp)
			consumer := client.New("localhost", cp)

			<-server.Started

			err := consumer.Subscribe(func(msg string) {
				incomingMessages.Write([]byte(msg))
			})
			assert.Nil(t, err)

			consumer2 := client.New("localhost", cp)
			err = consumer2.Subscribe(func(msg string) {
				incomingMessages.Write([]byte(msg))
			})
			assert.Nil(t, err)
			defer consumer2.Stop()

			publisher.SendMessage(message)
			publisher.SendMessage(message)

			consumer.Stop()

			publisher.SendMessage(message)
			publisher.SendMessage(message)
			publisher.SendMessage(message)

			assert.Eventually(t, func() bool {
				return assert.Contains(t, incomingMessages.Get(), message)
			}, time.Second, 10*time.Millisecond)

			assert.Eventually(t, func() bool {
				return assert.Len(t, incomingMessages.Get(), 3)
			}, time.Second, 10*time.Millisecond)
		})

		t.Run("it should log an error if publishers port is in use", func(t *testing.T) {
			sw := logger.SliceWriter{}
			l := log.New(&sw, "", log.Ldate)

			pp, cp := ":3339", ":3340"
			server := broker.NewBrokerServer(pp, cp, l)
			go server.Start()
			defer func() {
				go server.Stop()
			}()

			server2 := broker.NewBrokerServer(pp, ":3002", l)
			go server2.Start()

			assert.Eventually(t, func() bool {
				msg := time.Now().Format("2006/01/02") + " listen tcp " + pp + ": bind: address already in use"

				return assert.Contains(t, sw.Get(), msg)
			}, time.Second, 10*time.Millisecond)
		})

		t.Run("it should log an error if consumers port is in use", func(t *testing.T) {
			sw := logger.SliceWriter{}
			l := log.New(&sw, "", log.Ldate)

			pp, cp := ":3341", ":3342"
			server := broker.NewBrokerServer(pp, cp, l)
			go server.Start()
			defer func() {
				go server.Stop()
			}()

			server2 := broker.NewBrokerServer(":3002", cp, l)
			go server2.Start()

			assert.Eventually(t, func() bool {
				msg := time.Now().Format("2006/01/02") + " listen tcp " + cp + ": bind: address already in use"

				return assert.Contains(t, sw.Get(), msg)
			}, time.Second, 10*time.Millisecond)
		})

		t.Run("it should log an error if publishers port is invalid", func(t *testing.T) {
			sw := logger.SliceWriter{}
			l := log.New(&sw, "", log.Ldate)

			pp, cp := ":999999999", ":3343"
			server := broker.NewBrokerServer(pp, cp, l)
			go server.Start()

			assert.Eventually(t, func() bool {
				msg := time.Now().Format("2006/01/02") + " address " + strings.Trim(pp, ":") + ": invalid port"

				return assert.Contains(t, sw.Get(), msg)
			}, time.Second, 10*time.Millisecond)
		})

		t.Run("it should log an error if consumers port is invalid", func(t *testing.T) {
			sw := logger.SliceWriter{}
			l := log.New(&sw, "", log.Ldate)

			pp, cp := ":3344", ":99999999"
			server := broker.NewBrokerServer(pp, cp, l)
			go server.Start()

			assert.Eventually(t, func() bool {
				msg := time.Now().Format("2006/01/02") + " address " + strings.Trim(cp, ":") + ": invalid port"

				return assert.Contains(t, sw.Get(), msg)
			}, time.Second, 10*time.Millisecond)
		})
	})
}
