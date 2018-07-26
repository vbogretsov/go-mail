package mock

import (
	"container/list"
	"sync"

	"github.com/vbogretsov/go-mail"
)

// Sender represent mock for mailcd.Sender.
type Sender struct {
	mutex   sync.Mutex
	inboxes map[string]*list.List
	Error   error
}

// New creates new sender mock.
func New() *Sender {
	return &Sender{
		mutex:   sync.Mutex{},
		inboxes: map[string]*list.List{},
		Error:   nil,
	}
}

// Send implements mailcd.Sender.Send.
func (s *Sender) Send(req mailcd.Request) error {
	s.mutex.Lock()
	for _, addr := range req.To {
		s.send(req, addr.Email)
	}
	for _, addr := range req.Cc {
		s.send(req, addr.Email)
	}
	for _, addr := range req.Bcc {
		s.send(req, addr.Email)
	}
	s.mutex.Unlock()
	return s.Error
}

// Close implements mailcd.Sender.Close.
func (s *Sender) Close() error {
	return nil
}

func (s *Sender) ReadMail(email string) (mailcd.Request, bool) {
	var req mailcd.Request
	var ok bool
	s.mutex.Lock()
	inbox, ok := s.inboxes[email]
	ok = ok && inbox.Len() > 0
	if ok {
		node := inbox.Front()
		inbox.Remove(node)
		req = node.Value.(mailcd.Request)
	}
	s.mutex.Unlock()
	return req, ok
}

func (s *Sender) send(req mailcd.Request, recipient string) {
	if inbox, ok := s.inboxes[recipient]; ok {
		inbox.PushBack(req)
	} else {
		inbox := list.New()
		inbox.PushBack(req)
		s.inboxes[recipient] = inbox
	}
}
