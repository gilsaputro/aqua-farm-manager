package trackingevent

import (
	"aqua-farm-manager/internal/domain/stat"
	"regexp"
	"time"

	"encoding/json"
	"fmt"
	"strings"

	"github.com/nsqio/go-nsq"
)

// TrackingEventConsumer is list dependencies of consumer
type TrackingEventConsumer struct {
	topic        string
	channel      string
	host         string
	maxInFlight  int
	numconsumer  int
	timeoutInSec int
	stat         stat.StatDomain
}

// NewTrackingEverntConsumer is func to create TrackingEventConsumer
func NewTrackingEverntConsumer(topic, channel, host string, maxInFlight, numconsumer, timeoutInSec int, stat stat.StatDomain) *TrackingEventConsumer {
	return &TrackingEventConsumer{
		topic:        topic,
		channel:      channel,
		host:         host,
		stat:         stat,
		maxInFlight:  maxInFlight,
		numconsumer:  numconsumer,
		timeoutInSec: timeoutInSec,
	}
}

// Start is func to start the consumer
func (c *TrackingEventConsumer) Start() error {
	config := nsq.NewConfig()
	config.MaxInFlight = c.maxInFlight
	config.MsgTimeout = time.Duration(c.timeoutInSec) * time.Second
	consumer, err := nsq.NewConsumer(c.topic, c.channel, config)
	if err != nil {
		return err
	}

	consumer.AddConcurrentHandlers(c, c.numconsumer)

	lookupd := strings.Split(c.host, ",")
	if len(lookupd) < 1 {
		return fmt.Errorf("invalid lookupd config")
	}
	return consumer.ConnectToNSQLookupds(lookupd)
}

// HandleMessage is func to handler the message from aqua_farm_tracking_event
func (c *TrackingEventConsumer) HandleMessage(msg *nsq.Message) error {
	var err error
	msg.DisableAutoResponse()
	var body TrackingEventMessage
	defer msg.Finish()
	err = json.Unmarshal(msg.Body, &body)
	if err != nil {
		fmt.Println("TrackingEventConsumer-Got Error Unmarshal :", err)
		return nil
	}

	// checking valid body
	if len(body.Path) < 1 || len(body.Method) < 1 || len(body.UA) < 1 {
		return nil
	}

	// Validate path contain ID
	path := body.Path
	// Should be change to apply global path
	re := regexp.MustCompile(`^/(farms|ponds)/(\d+)$`)
	split := re.FindStringSubmatch(path)
	if len(split) > 1 {
		path = "/" + split[1]
	}

	c.stat.IngestStatAPI(stat.IngestStatRequest{
		Path:   path,
		Method: body.Method,
		Ua:     body.UA,
		Code:   body.Code,
	})
	return nil
}
