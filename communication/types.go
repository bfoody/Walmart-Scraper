package communication

import "github.com/bfoody/Walmart-Scraper/domain"

// A SingleReceiverPacket is a message meant to be received by a single client.
type SingleReceiverPacket struct {
	SenderID   string
	ReceiverID string
}

// A FanoutPacket is a message meant to be received by all clients.
type FanoutPacket struct {
	SenderID string
}

// A Heartbeat is sent to another server to notify it that the sending server is still healthy.
type Heartbeat struct {
	SingleReceiverPacket
	ResponseExpected bool // whether or not a heartbeat is required in response
}

// A StatusUpdate is sent by a server to notify others of an update of capabilities.
type StatusUpdate struct {
	FanoutPacket
	AvailableForWork bool // whether or not the server can be assigned work yet
}

// A HubWelcome is sent to a client when the hub registers it.
type HubWelcome struct {
	SingleReceiverPacket
}

// A HubWelcomeAck is sent to the hub from the client to indicate
// that the receiving hub is now the client's primary hub.
type HubWelcomeAck struct {
	SingleReceiverPacket
}

// A GoingAway is sent by a client to a hub to let it know that it will
// be shutting down and to gracefully remove it from its supervisor.
type GoingAway struct {
	SingleReceiverPacket
	Reason string
}

// An InfoRetrieved is sent by a client after fetching product info for a hub. The `ID` field in the ProductInfo will be blank.
type InfoRetrieved struct {
	SingleReceiverPacket
	TaskID      string
	ProductInfo domain.ProductInfo
}

// A TaskFTaskFulfillmentRequest is sent by a hub to a client as a request for a task to be executed and fulfilled.
type TaskFulfillmentRequest struct {
	SingleReceiverPacket
	TaskID          string
	ProductLocation domain.ProductLocation
}
