package identity

const (
	// VariantHub represents a hub or parent server.
	VariantHub = "hub"
	// VariantClient represents a client or child server.
	VariantClient = "client"
)

// A Server represents a server's identity, including its ID and properties.
type Server struct {
	Variant string // can be `hub` or `client`
	ID      string
}

// NewServer creates and returns a *Server with the provided server ID.
func NewServer(variant string, id string) *Server {
	return &Server{
		Variant: variant,
		ID:      id,
	}
}

// NewHub is shorthand for NewServer(VariantHub, id).
func NewHub(id string) *Server {
	return NewServer(VariantHub, id)
}

// NewClient is shorthand for NewServer(VariantHub, id).
func NewClient(id string) *Server {
	return NewServer(VariantClient, id)
}
