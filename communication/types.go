package communication

// A SingleReceiverPacket is a message meant to be received by a single client.
type SingleReceiverPacket struct {
	SenderID   string // the ID of the server sending the message
	ReceiverID string // the ID of the server receiving the message
}

// A Heartbeat is sent to another server to notify it that the sending server is still healthy.
type Heartbeat struct {
	SingleReceiverPacket
	ResponseExpected bool // whether or not a heartbeat is required in response
}
