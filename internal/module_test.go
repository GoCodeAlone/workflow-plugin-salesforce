package internal

import (
	"context"
	"testing"
)

func TestModuleInit_WithAccessToken(t *testing.T) {
	m, err := newSalesforceModule("test-token", map[string]any{
		"accessToken": "test-access-token",
		"instanceUrl": "https://myorg.salesforce.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Init(); err != nil {
		t.Fatal(err)
	}
	c, ok := GetClient("test-token")
	if !ok || c == nil {
		t.Error("expected client to be registered")
	}
	UnregisterClient("test-token")
}

func TestModuleStop_UnregistersClient(t *testing.T) {
	m, _ := newSalesforceModule("test-stop", map[string]any{
		"accessToken": "tok",
		"instanceUrl": "https://myorg.salesforce.com",
	})
	_ = m.Init()
	_ = m.Stop(context.Background())
	_, ok := GetClient("test-stop")
	if ok {
		t.Error("expected client to be unregistered after stop")
	}
}

func TestModuleInit_MissingCredentials(t *testing.T) {
	m, err := newSalesforceModule("test-missing", map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Init(); err == nil {
		t.Error("expected error for missing credentials")
		UnregisterClient("test-missing")
	}
}

func TestModuleInit_AccessTokenMissingInstanceUrl(t *testing.T) {
	m, err := newSalesforceModule("test-noinstance", map[string]any{
		"accessToken": "tok",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Init(); err == nil {
		t.Error("expected error for missing instanceUrl")
		UnregisterClient("test-noinstance")
	}
}

func TestModuleInit_DefaultAPIVersion(t *testing.T) {
	m, _ := newSalesforceModule("test-version", map[string]any{
		"accessToken": "tok",
		"instanceUrl": "https://myorg.salesforce.com",
	})
	_ = m.Init()
	c, ok := GetClient("test-version")
	if !ok {
		t.Fatal("expected client")
	}
	if c.apiVersion != defaultAPIVersion {
		t.Errorf("expected apiVersion %q, got %q", defaultAPIVersion, c.apiVersion)
	}
	UnregisterClient("test-version")
}

func TestModuleInit_CustomAPIVersion(t *testing.T) {
	m, _ := newSalesforceModule("test-custom-version", map[string]any{
		"accessToken": "tok",
		"instanceUrl": "https://myorg.salesforce.com",
		"apiVersion":  "v60.0",
	})
	_ = m.Init()
	c, ok := GetClient("test-custom-version")
	if !ok {
		t.Fatal("expected client")
	}
	if c.apiVersion != "v60.0" {
		t.Errorf("expected apiVersion %q, got %q", "v60.0", c.apiVersion)
	}
	UnregisterClient("test-custom-version")
}
