package nsq

import (
	"encoding/json"

	"github.com/nsqio/go-nsq"
)

// NsqMethod is list all available method for nsq
type NsqMethod interface {
	Publish(topic string, data interface{}) error
}

// Client is a wrapper for Postgres client
type Client struct {
	nsq *nsq.Producer
}

// NewNsqClient is func to create nsq client
func NewNsqClient(host string) (NsqMethod, error) {
	conf := nsq.NewConfig()
	producer, err := nsq.NewProducer(host, conf)
	if err != nil {
		return &Client{}, err
	}

	return &Client{nsq: producer}, nil
}

// Publish is func to publisn message
func (c *Client) Publish(topic string, data interface{}) error {
	var body []byte

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return c.nsq.Publish(topic, body)
}
