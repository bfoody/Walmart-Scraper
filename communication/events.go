package communication

import (
	"encoding/json"
	"strings"
)

// A QueueConnection wraps an AMQP connection and allows for event handlers to be registered.
type QueueConnection struct {
	conn                           *Connection
	queueName                      string
	heartbeatHandler               func(heartbeat *Heartbeat)
	statusUpdateHandler            func(statusUpdate *StatusUpdate)
	hubWelcomeHandler              func(hubWelcome *HubWelcome)
	hubWelcomeAckHandler           func(hubWelcomeAck *HubWelcomeAck)
	goingAwayHandler               func(goingAway *GoingAway)
	infoRetrievedHandler           func(infoRetrieved *InfoRetrieved)
	taskFulfillmentRequestHandler  func(taskFulfillmentRequest *TaskFulfillmentRequest)
	crawlFulfillmentRequestHandler func(crawlFulfillmentRequest *CrawlFulfillmentRequest)
	crawlRetrievedHandler          func(crawlRetrieved *CrawlRetrieved)
}

// NewQueueConnection creates and returns a new QueueConnection.
func NewQueueConnection(conn *Connection, queueName string) *QueueConnection {
	return &QueueConnection{
		conn:      conn,
		queueName: queueName,
	}
}

// Consume starts consuming from the AMQP channel.
func (q *QueueConnection) Consume() error {
	channel, err := q.conn.Subscribe(q.queueName)
	if err != nil {
		return err
	}

	go q.consumer(channel)

	return nil
}

// consumer consumes messages from the queue.
func (q *QueueConnection) consumer(channel chan Message) {
	for {
		msg := <-channel

		decoder := json.NewDecoder(strings.NewReader(string(msg.Content)))

		switch msg.Type {
		case "heartbeat":
			d := &Heartbeat{}
			if err := decoder.Decode(d); err == nil && q.heartbeatHandler != nil {
				q.heartbeatHandler(d)
			}
			break
		case "statusUpdate":
			d := &StatusUpdate{}
			if err := decoder.Decode(d); err == nil && q.statusUpdateHandler != nil {
				q.statusUpdateHandler(d)
			}
			break
		case "hubWelcome":
			d := &HubWelcome{}
			if err := decoder.Decode(d); err == nil && q.hubWelcomeHandler != nil {
				q.hubWelcomeHandler(d)
			}
			break
		case "hubWelcomeAck":
			d := &HubWelcomeAck{}
			if err := decoder.Decode(d); err == nil && q.hubWelcomeAckHandler != nil {
				q.hubWelcomeAckHandler(d)
			}
			break
		case "goingAway":
			d := &GoingAway{}
			if err := decoder.Decode(d); err == nil && q.goingAwayHandler != nil {
				q.goingAwayHandler(d)
			}
			break
		case "infoRetrieved":
			d := &InfoRetrieved{}
			if err := decoder.Decode(d); err == nil && q.infoRetrievedHandler != nil {
				q.infoRetrievedHandler(d)
			}
			break
		case "taskFulfillmentRequest":
			d := &TaskFulfillmentRequest{}
			if err := decoder.Decode(d); err == nil && q.taskFulfillmentRequestHandler != nil {
				q.taskFulfillmentRequestHandler(d)
			}
			break
		case "crawlFulfillmentRequest":
			d := &CrawlFulfillmentRequest{}
			if err := decoder.Decode(d); err == nil && q.crawlFulfillmentRequestHandler != nil {
				q.crawlFulfillmentRequestHandler(d)
			}
			break
		case "crawlRetrieved":
			d := &CrawlRetrieved{}
			if err := decoder.Decode(d); err == nil && q.crawlRetrievedHandler != nil {
				q.crawlRetrievedHandler(d)
			}
			break
		}
	}
}

// RegisterHeartbeatHandler registers a handler for Heartbeat messages.
func (q *QueueConnection) RegisterHeartbeatHandler(handler func(heartbeat *Heartbeat)) {
	q.heartbeatHandler = handler
}

// RegisterStatusUpdateHandler registers a handler for StatusUpdate messages.
func (q *QueueConnection) RegisterStatusUpdateHandler(handler func(statusUpdate *StatusUpdate)) {
	q.statusUpdateHandler = handler
}

// RegisterHubWelcomeHandler registers a handler for HubWelcome messages.
func (q *QueueConnection) RegisterHubWelcomeHandler(handler func(hubWelcome *HubWelcome)) {
	q.hubWelcomeHandler = handler
}

// RegisterHubWelcomeAckHandler registers a handler for HubWelcomeAck messages.
func (q *QueueConnection) RegisterHubWelcomeAckHandler(handler func(hubWelcomeAck *HubWelcomeAck)) {
	q.hubWelcomeAckHandler = handler
}

// RegisterGoingAwayHandler registers a handler for GoingAway messages.
func (q *QueueConnection) RegisterGoingAwayHandler(handler func(goingAway *GoingAway)) {
	q.goingAwayHandler = handler
}

// RegisterInfoRetrievedHandler registers a handler for InfoRetrieved messages.
func (q *QueueConnection) RegisterInfoRetrievedHandler(handler func(infoRetrieved *InfoRetrieved)) {
	q.infoRetrievedHandler = handler
}

// RegisterTaskFulfillmentRequest registers a handler for TaskFulfillmentRequest messages.
func (q *QueueConnection) RegisterTaskFulfillmentRequest(handler func(taskFulfillmentRequest *TaskFulfillmentRequest)) {
	q.taskFulfillmentRequestHandler = handler
}

// RegisterCrawlFulfillmentRequest registers a handler for CrawlFulfillmentRequest messages.
func (q *QueueConnection) RegisterCrawlFulfillmentRequestHandler(handler func(crawlFulfillmentRequest *CrawlFulfillmentRequest)) {
	q.crawlFulfillmentRequestHandler = handler
}

// RegisterCrawlRetrieved registers a handler for CrawlRetrieved messages.
func (q *QueueConnection) RegisterCrawlRetrievedHandler(handler func(crawlRetrieved *CrawlRetrieved)) {
	q.crawlRetrievedHandler = handler
}

// SendMessage sends a message of any supported type to the queue,
// panicking if an invalid type is sent.
func (q *QueueConnection) SendMessage(message interface{}) error {
	typeName := ""

	switch message.(type) {
	case Heartbeat:
		typeName = "heartbeat"
	case StatusUpdate:
		typeName = "statusUpdate"
	case HubWelcome:
		typeName = "hubWelcome"
	case HubWelcomeAck:
		typeName = "hubWelcomeAck"
	case GoingAway:
		typeName = "goingAway"
	case InfoRetrieved:
		typeName = "infoRetrieved"
	case TaskFulfillmentRequest:
		typeName = "taskFulfillmentRequest"
	case CrawlFulfillmentRequest:
		typeName = "crawlFulfillmentRequest"
	case CrawlRetrieved:
		typeName = "crawlRetrieved"
	}

	if typeName != "" {
		return q.conn.Send(q.queueName, typeName, message)
	}

	panic("invalid message type")
}
