package internal

import (
	"sync"
)

var (
	clientMu       sync.RWMutex
	clientRegistry = make(map[string]*salesforceClient)
)

// RegisterClient adds a Salesforce client to the global registry under the given name.
func RegisterClient(name string, c *salesforceClient) {
	clientMu.Lock()
	defer clientMu.Unlock()
	clientRegistry[name] = c
}

// GetClient looks up a Salesforce client by name.
func GetClient(name string) (*salesforceClient, bool) {
	clientMu.RLock()
	defer clientMu.RUnlock()
	c, ok := clientRegistry[name]
	return c, ok
}

// UnregisterClient removes a client from the registry.
func UnregisterClient(name string) {
	clientMu.Lock()
	defer clientMu.Unlock()
	delete(clientRegistry, name)
}
