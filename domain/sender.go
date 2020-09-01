package domain

import (
	"strconv"
	"sync"
)

type Sender struct {
	clients     map[Client]struct{}
	clientMutex *sync.RWMutex
}

func NewSender() *Sender {
	return &Sender{
		clients:     make(map[Client]struct{}),
		clientMutex: &sync.RWMutex{},
	}
}

func (s *Sender) AddClient(c Client) error {
	s.clientMutex.Lock()
	defer s.clientMutex.Unlock()
	s.clients[c] = struct{}{}
	return nil
}

func (s *Sender) DeleteClient(c Client) {
	s.clientMutex.Lock()
	defer s.clientMutex.Unlock()
	delete(s.clients, c)
	c.Close()
}

func (s *Sender) SendMessage(m Message) {
	s.clientMutex.RLock()
	defer s.clientMutex.RUnlock()
	for client := range s.clients {
		client.SendMessage(m)
	}
}

func (s *Sender) Stats() string {
	return "Client count = " + strconv.Itoa(len(s.clients))
}
