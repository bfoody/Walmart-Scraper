package communication

import (
	"encoding/json"
	"strings"

	"github.com/streadway/amqp"
)

// A Connection represents a connection to the message queue through AMQP.
type Connection struct {
	conn *amqp.Connection
}

// A Message represents a message sent through the AMQP message queue.
type Message struct {
	Type    string
	Content json.RawMessage // should be parsed based on `Type`
}

// ConnectAMQP attempts to dial an AMQP server, returning a *Connection on
// success.
func ConnectAMQP(url string) (*Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	return &Connection{
		conn,
	}, nil
}

// Subscribe subscribes to a queue, returning a channel of Messages.
func (c *Connection) Subscribe(queue string) (chan Message, error) {
	ch, err := c.conn.Channel()
	if err != nil {
		return nil, err
	}

	channel := make(chan Message, 2)

	in, err := ch.Consume(queue, "", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	go func() {
		for message := range in {
			body := string(message.Body)

			var msg Message
			err := json.NewDecoder(strings.NewReader(body)).Decode(&msg)
			if err != nil {
				// TODO: Find a better way to deal with this error.
				continue
			}

			channel <- msg
		}
	}()

	return channel, nil
}
