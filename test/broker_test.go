package test

import (
	"github.com/andriusbil/tcp-broker/broker"
	"github.com/andriusbil/tcp-broker/client"
	"github.com/andriusbil/tcp-broker/logger"
	"github.com/stretchr/testify/assert"
	"log"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestBroker(t *testing.T) {
	t.Run("Broker server", func(t *testing.T) {
		t.Run("it should notify consumers with new messages from publishers", func(t *testing.T) {
			message := "Test123"

			type Messages struct {
				mu    sync.Mutex
				Value []string
			}

			incomingMessages := Messages{}

			pp, cp := ":3000", ":3001"
			server := broker.NewBrokerServer(pp, cp, log.Default())
			go server.Start()
			defer server.Stop()

			publisher := client.New("localhost", pp)
			publisher2 := client.New("localhost", pp)
			consumer := client.New("localhost", cp)
			consumer2 := client.New("localhost", cp)

			consumer.Subscribe(func(msg string) {
				incomingMessages.mu.Lock()
				incomingMessages.Value = append(incomingMessages.Value, "node1:"+msg)
				incomingMessages.mu.Unlock()
			})
			defer consumer.Stop()

			consumer2.Subscribe(func(msg string) {
				incomingMessages.mu.Lock()
				incomingMessages.Value = append(incomingMessages.Value, "node2:"+msg)
				incomingMessages.mu.Unlock()
			})
			defer consumer2.Stop()

			publisher.SendMessage(message)
			publisher.SendMessage(message)
			publisher2.SendMessage(message)
			publisher2.SendMessage(message)

			assert.Eventually(t, func() bool {
				return assert.Contains(t, incomingMessages.Value, "node1:"+message)
			}, time.Second, 10*time.Millisecond)

			assert.Eventually(t, func() bool {
				return assert.Contains(t, incomingMessages.Value, "node2:"+message)
			}, time.Second, 10*time.Millisecond)
		})

		t.Run("it should notify new consumers with previously published messages", func(t *testing.T) {
			message := "Test123"

			type Messages struct {
				mu    sync.Mutex
				Value []string
			}

			incomingMessages := Messages{}

			pp, cp := ":3000", ":3001"
			server := broker.NewBrokerServer(pp, cp, log.Default())
			go server.Start()
			defer server.Stop()

			publisher := client.New("localhost", pp)
			consumer := client.New("localhost", cp)

			publisher.SendMessage(message)
			publisher.SendMessage(message)
			publisher.SendMessage(message)
			publisher.SendMessage(message)

			// Subscribe after push
			consumer.Subscribe(func(msg string) {
				incomingMessages.mu.Lock()
				incomingMessages.Value = append(incomingMessages.Value, "node1:"+msg)
				incomingMessages.mu.Unlock()
			})
			defer consumer.Stop()

			assert.Eventually(t, func() bool {
				return assert.Contains(t, incomingMessages.Value, "node1:"+message)
			}, time.Second, 10*time.Millisecond)

			assert.Eventually(t, func() bool {
				return assert.Len(t, incomingMessages.Value, 4)
			}, time.Second, 10*time.Millisecond)
		})

		t.Run("it should continue serving consumers if some of them stop", func(t *testing.T) {
			message := "Test123"

			type Messages struct {
				mu    sync.Mutex
				Value []string
			}

			incomingMessages := Messages{}

			pp, cp := ":3000", ":3001"
			server := broker.NewBrokerServer(pp, cp, log.Default())
			go server.Start()
			defer server.Stop()

			publisher := client.New("localhost", pp)
			consumer := client.New("localhost", cp)
			consumer.Subscribe(func(msg string) {
				incomingMessages.mu.Lock()
				incomingMessages.Value = append(incomingMessages.Value, "node1:"+msg)
				incomingMessages.mu.Unlock()
			})

			consumer2 := client.New("localhost", cp)
			consumer2.Subscribe(func(msg string) {
				incomingMessages.mu.Lock()
				incomingMessages.Value = append(incomingMessages.Value, "node2:"+msg)
				incomingMessages.mu.Unlock()
			})
			defer consumer2.Stop()

			publisher.SendMessage(message)
			publisher.SendMessage(message)

			consumer.Stop()

			publisher.SendMessage(message)
			publisher.SendMessage(message)
			publisher.SendMessage(message)

			assert.Eventually(t, func() bool {
				return assert.Contains(t, incomingMessages.Value, "node1:"+message)
			}, time.Second, 10*time.Millisecond)

			assert.Eventually(t, func() bool {
				return assert.Contains(t, incomingMessages.Value, "node2:"+message)
			}, time.Second, 10*time.Millisecond)

			assert.Eventually(t, func() bool {
				return len(incomingMessages.Value) >= 3
			}, time.Second, 10*time.Millisecond)
		})

		t.Run("it should log an error if port is in use", func(t *testing.T) {
			sw := logger.SliceWriter{}
			l := log.New(&sw, "", log.Ldate)

			pp, cp := ":3000", ":3001"
			server := broker.NewBrokerServer(pp, cp, l)
			go server.Start()
			defer func() {
				go server.Stop()
			}()

			server2 := broker.NewBrokerServer(pp, cp, l)
			go server2.Start()
			defer func() {
				go server2.Stop()
			}()

			assert.Eventually(t, func() bool {
				msg := time.Now().Format("2006/01/02") + " listen tcp " + pp + ": bind: address already in use"
				msg2 := time.Now().Format("2006/01/02") + " listen tcp " + cp + ": bind: address already in use"

				for _, s := range sw.Slice {
					if s == msg || s == msg2 {
						return true
					}
				}

				return false
			}, time.Second, 10*time.Millisecond)
		})

		t.Run("it should log an error if port is invalid", func(t *testing.T) {
			sw := logger.SliceWriter{}
			l := log.New(&sw, "", log.Ldate)

			pp, cp := ":999999999", ":999999999"
			server := broker.NewBrokerServer(pp, cp, l)
			go server.Start()
			defer func() {
				go server.Stop()
			}()

			assert.Eventually(t, func() bool {
				msg := time.Now().Format("2006/01/02") + " address " + strings.Trim(pp, ":") + ": invalid port"

				return assert.Contains(t, sw.Slice, msg)
			}, time.Second, 10*time.Millisecond)
		})
	})
}
