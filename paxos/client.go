package paxos

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	RecvTimeoutError = errors.New("recv timeout")
)

type Client interface {
	GetId() int
	Send(msg Message) error
	Recv(timeout time.Duration) (Message, error)
}

func NewEagerNetwork() *EagerNetwork {
	return &EagerNetwork{new(sync.RWMutex), make(map[int]chan Message)}
}

type EagerNetwork struct {
	*sync.RWMutex
	recvQueue map[int]chan Message
}

func (en *EagerNetwork) Join(clientId int) {
	en.Lock()
	defer en.Unlock()
	en.recvQueue[clientId] = make(chan Message)
}

type EagerClient struct {
	clientId int
	EagerNetwork
}

func (client *EagerClient) GetId() int {
	return client.clientId
}

func (client *EagerClient) Send(msg Message) error {
	ch, exists := client.recvQueue[msg.GetTo()]
	if !exists {
		return fmt.Errorf("failed to send msg. unavailable address is: %d", msg.GetTo())
	}
	ch <- msg
	return nil
}

func (client *EagerClient) Recv(timeout time.Duration) (Message, error) {
	ch, exists := client.recvQueue[client.clientId]
	if !exists {
		return Message{}, fmt.Errorf("failed to recv msg. unavailable address is: %d", client.clientId)
	}

	select {
	case <-time.After(timeout):
		return Message{}, RecvTimeoutError
	case msg := <-ch:
		return msg, nil
	}
}
