package communication

import (
	"encoding/json"
	"strings"
)

// A QueueConnection wraps an AMQP connection and allows for event handlers to be registered.
type QueueConnection struct {
	conn                *Connection
	queueName           string
	heartbeatHandler    func(heartbeat *Heartbeat)
	statusUpdateHandler func(statusUpdate *StatusUpdate)
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
		case "heatbeat":
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

func (q *QueueConnection) SendMessage(message interface{}) error {
	switch message.(type) {
	case Heartbeat:
		return q.conn.Send(q.queueName, "heartbeat", message)
	case StatusUpdate:
		return q.conn.Send(q.queueName, "statusUpdate", message)
	}

	panic("invalid message type")
}
