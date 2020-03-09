package micro

// Client is a wrapper interface for Micro clients.
type Client interface {
	Name() string
}

// ClientSet is a container for multiple clients used by a Service.
type ClientSet struct {
	clients map[string]Client
}

// NewClientSet creates a new ClientSet with the given clients.
func NewClientSet(clients ...Client) *ClientSet {
	c := &ClientSet{make(map[string]Client)}
	c.AddClient(clients...)

	return c
}

// AddClient adds a client to a ClientSet.
func (c *ClientSet) AddClient(clients ...Client) {
	for _, client := range clients {
		c.clients[client.Name()] = client
	}
}

// Get returns the Client for a given Service name.
func (c *ClientSet) Get(key string) Client {
	return c.clients[key]
}

// Keys returns the service names of clients in the ClientSet.
func (c *ClientSet) Keys() []string {
	keys := []string{}

	for key := range c.clients {
		keys = append(keys, key)
	}

	return keys
}

// Length returns the number of clients in the ClientSet.
func (c *ClientSet) Length() int {
	return len(c.clients)
}
