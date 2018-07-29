package mock

import (
	"container/list"
	"sync"

	"github.com/vbogretsov/go-mail"
)

type mailbox struct {
	mutex sync.Mutex
	items *list.List
}

// Sender represent mock for mail.Sender.
type Sender struct {
	mutex sync.RWMutex
	boxes map[string]*mailbox
	Error error
}

// New creates new sender mock.
func New() *Sender {
	return &Sender{
		boxes: map[string]*mailbox{},
		Error: nil,
	}
}

// Send implements mail.Sender.Send.
func (s *Sender) Send(req mail.Request) error {
	for _, addr := range req.To {
		s.send(req, addr.Email)
	}
	for _, addr := range req.Cc {
		s.send(req, addr.Email)
	}
	for _, addr := range req.Bcc {
		s.send(req, addr.Email)
	}
	return s.Error
}

// Close implements mail.Sender.Close.
func (s *Sender) Close() error {
	return nil
}

func (s *Sender) ReadMail(email string) (mail.Request, bool) {
	var req mail.Request
	var ok bool

	s.mutex.RLock()
	box, ok := s.boxes[email]
	s.mutex.RUnlock()

	if ok {
		box.mutex.Lock()
		if box.items.Len() > 0 {
			node := box.items.Front()
			box.items.Remove(node)
			req = node.Value.(mail.Request)
		} else {
			ok = false
		}
		box.mutex.Unlock()
	}

	return req, ok
}

func (s *Sender) send(req mail.Request, recipient string) {
	s.mutex.RLock()
	box, ok := s.boxes[recipient]
	s.mutex.RUnlock()

	if !ok {
		box = &mailbox{items: list.New()}
		s.mutex.Lock()
		s.boxes[recipient] = box
		s.mutex.Unlock()
	}

	box.mutex.Lock()
	box.items.PushBack(req)
	box.mutex.Unlock()
}
