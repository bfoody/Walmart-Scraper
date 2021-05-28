package communication

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
