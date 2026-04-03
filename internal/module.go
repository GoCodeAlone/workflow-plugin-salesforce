package internal

import (
	"context"
	"fmt"

	"github.com/GoCodeAlone/workflow-plugin-salesforce/salesforce"
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

	// Otherwise use OAuth client credentials via the provider package
	if clientID == "" || clientSecret == "" {
		return fmt.Errorf("salesforce.provider %q: clientId and clientSecret are required (or provide accessToken + instanceUrl)", m.name)
	}

	cfg := salesforce.Config{
		AuthType:     "client_credentials",
		ClientID:     clientID,
		ClientSecret: clientSecret,
		LoginURL:     loginURL,
		InstanceURL:  instanceURL,
		APIVersion:   apiVersion,
	}

	provider, err := salesforce.NewProvider(context.Background(), cfg)
	if err != nil {
		return fmt.Errorf("salesforce.provider %q: auth failed: %w", m.name, err)
	}

	if apiVersion == "" {
		apiVersion = defaultAPIVersion
	}

	client := newSalesforceClientFromSDK(provider.Client, apiVersion)
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
