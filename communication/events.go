package communication

import (
	"encoding/json"
	"fmt"
	"strings"
)

// A QueueConnection wraps an AMQP connection and allows for event handlers to be registered.
type QueueConnection struct {
	conn                 *Connection
	queueName            string
	heartbeatHandler     func(heartbeat *Heartbeat)
	statusUpdateHandler  func(statusUpdate *StatusUpdate)
	hubWelcomeHandler    func(hubWelcome *HubWelcome)
	hubWelcomeAckHandler func(hubWelcomeAck *HubWelcomeAck)
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

		fmt.Println(msg)

		decoder := json.NewDecoder(strings.NewReader(string(msg.Content)))

		switch msg.Type {
		case "heartbeat":
			fmt.Println("heartbeat")
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
	}

	if typeName != "" {
		return q.conn.Send(q.queueName, typeName, message)
	}

	panic("invalid message type")
}
