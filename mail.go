package mailcd

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

// Address represents an email address.
type Address struct {
	Email string `json:"email,omitempty"`
	Name  string `json:"name,omitempty"`
}

// Request represents a maild request.
type Request struct {
	TemplateLang string                 `json:"templateLang"`
	TemplateName string                 `json:"templateName"`
	TemplateArgs map[string]interface{} `json:"templateArgs"`
	To           []Address              `json:"to"`
	Cc           []Address              `json:"cc"`
	Bcc          []Address              `json:"bcc"`
}

// Sender represents interface of an email sender.
type Sender interface {
	Send(Request) error
	Close() error
}

type amqpSender struct {
	channel *amqp.Channel
	topic   string
}

func (s *amqpSender) Send(req Request) error {
	buf, err := json.Marshal(req)
	if err != nil {
		return err
	}

	msg := amqp.Publishing{Body: buf}

	return s.channel.Publish(s.topic, s.topic, false, false, msg)
}

func (s *amqpSender) Close() error {
	return s.channel.Close()
}

// New creates a new maild client.
func NewSender(conn *amqp.Connection, topic string) (Sender, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	sender := amqpSender{
		channel: ch,
		topic:   topic,
	}
	return &sender, nil
}
