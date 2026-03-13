package internal

import (
	"context"
	"fmt"
)

// salesforceModule creates a Salesforce REST client and registers it.
type salesforceModule struct {
	name   string
	config map[string]any
}

func newSalesforceModule(name string, config map[string]any) (*salesforceModule, error) {
	return &salesforceModule{name: name, config: config}, nil
}

// Init creates the Salesforce REST client and registers it in the global registry.
func (m *salesforceModule) Init() error {
	loginURL, _ := m.config["loginUrl"].(string)
	clientID, _ := m.config["clientId"].(string)
	clientSecret, _ := m.config["clientSecret"].(string)
	accessToken, _ := m.config["accessToken"].(string)
	instanceURL, _ := m.config["instanceUrl"].(string)
	apiVersion, _ := m.config["apiVersion"].(string)

	// If direct access token provided, use it
	if accessToken != "" {
		if instanceURL == "" {
			return fmt.Errorf("salesforce.provider %q: instanceUrl is required when using accessToken", m.name)
		}
		client := newSalesforceClient(instanceURL, accessToken, apiVersion)
		RegisterClient(m.name, client)
		return nil
	}

	// Otherwise use OAuth client credentials
	if clientID == "" || clientSecret == "" {
		return fmt.Errorf("salesforce.provider %q: clientId and clientSecret are required (or provide accessToken + instanceUrl)", m.name)
	}
	if loginURL == "" {
		loginURL = "https://login.salesforce.com"
	}

	iURL, token, err := authenticateOAuth(loginURL, clientID, clientSecret)
	if err != nil {
		return fmt.Errorf("salesforce.provider %q: auth failed: %w", m.name, err)
	}
	if instanceURL == "" {
		instanceURL = iURL
	}

	client := newSalesforceClient(instanceURL, token, apiVersion)
	RegisterClient(m.name, client)
	return nil
}

// Start is a no-op for this module.
func (m *salesforceModule) Start(_ context.Context) error { return nil }

// Stop unregisters the Salesforce client.
func (m *salesforceModule) Stop(_ context.Context) error {
	UnregisterClient(m.name)
	return nil
}
