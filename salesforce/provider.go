// Package salesforce provides an exported Salesforce client provider
// backed by the PramithaMJ Salesforce Go SDK v2. Other plugins (e.g.
// workflow-plugin-crm) can import this package to interact with
// Salesforce without duplicating authentication or HTTP logic.
package salesforce

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	sf "github.com/PramithaMJ/salesforce/v2"
	"github.com/PramithaMJ/salesforce/v2/types"
)

// Provider wraps the PramithaMJ Salesforce SDK client.
type Provider struct {
	Client *sf.Client
}

// Config holds the configuration for creating a Provider.
type Config struct {
	AuthType      string // "oauth_refresh", "password", "client_credentials", "access_token" (or empty to auto-detect)
	ClientID      string
	ClientSecret  string
	RefreshToken  string
	Username      string
	Password      string
	SecurityToken string
	AccessToken   string
	InstanceURL   string
	LoginURL      string
	APIVersion    string
	Sandbox       bool
}

// NewProvider creates and connects a Salesforce SDK client from the
// given configuration. The returned Provider exposes the fully-
// initialised *sf.Client for typed API operations.
func NewProvider(ctx context.Context, cfg Config) (*Provider, error) {
	var opts []sf.Option

	if cfg.APIVersion != "" {
		opts = append(opts, sf.WithAPIVersion(strings.TrimPrefix(cfg.APIVersion, "v")))
	}
	if cfg.Sandbox {
		opts = append(opts, sf.WithSandbox())
	}

	tokenURL := tokenEndpoint(cfg.LoginURL)
	needsConnect := false

	switch cfg.AuthType {
	case "oauth_refresh":
		opts = append(opts, sf.WithOAuthRefresh(cfg.ClientID, cfg.ClientSecret, cfg.RefreshToken))
		if tokenURL != "" {
			opts = append(opts, sf.WithTokenURL(tokenURL))
		}
		needsConnect = true

	case "password":
		opts = append(opts, sf.WithPasswordAuth(cfg.Username, cfg.Password, cfg.SecurityToken))
		opts = append(opts, sf.WithCredentials(cfg.ClientID, cfg.ClientSecret))
		if tokenURL != "" {
			opts = append(opts, sf.WithTokenURL(tokenURL))
		}
		needsConnect = true

	case "client_credentials":
		instanceURL, accessToken, err := AuthenticateClientCredentials(cfg.LoginURL, cfg.ClientID, cfg.ClientSecret)
		if err != nil {
			return nil, err
		}
		if cfg.InstanceURL == "" {
			cfg.InstanceURL = instanceURL
		}
		opts = append(opts, sf.WithAccessToken(accessToken, cfg.InstanceURL))

	case "access_token", "":
		if cfg.AccessToken != "" {
			opts = append(opts, sf.WithAccessToken(cfg.AccessToken, cfg.InstanceURL))
		} else if cfg.RefreshToken != "" {
			opts = append(opts, sf.WithOAuthRefresh(cfg.ClientID, cfg.ClientSecret, cfg.RefreshToken))
			if tokenURL != "" {
				opts = append(opts, sf.WithTokenURL(tokenURL))
			}
			needsConnect = true
		} else if cfg.Username != "" {
			opts = append(opts, sf.WithPasswordAuth(cfg.Username, cfg.Password, cfg.SecurityToken))
			opts = append(opts, sf.WithCredentials(cfg.ClientID, cfg.ClientSecret))
			if tokenURL != "" {
				opts = append(opts, sf.WithTokenURL(tokenURL))
			}
			needsConnect = true
		} else if cfg.ClientID != "" && cfg.ClientSecret != "" {
			instanceURL, accessToken, err := AuthenticateClientCredentials(cfg.LoginURL, cfg.ClientID, cfg.ClientSecret)
			if err != nil {
				return nil, err
			}
			if cfg.InstanceURL == "" {
				cfg.InstanceURL = instanceURL
			}
			opts = append(opts, sf.WithAccessToken(accessToken, cfg.InstanceURL))
		} else {
			return nil, fmt.Errorf("salesforce: no valid authentication configuration provided")
		}

	default:
		return nil, fmt.Errorf("salesforce: unsupported auth type %q", cfg.AuthType)
	}

	client, err := sf.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("salesforce: create client: %w", err)
	}

	if needsConnect {
		if err := client.Connect(ctx); err != nil {
			return nil, fmt.Errorf("salesforce: connect: %w", err)
		}
	}

	return &Provider{Client: client}, nil
}

// AuthenticateClientCredentials performs the OAuth 2.0 client_credentials
// grant flow and returns the instance URL and access token.
func AuthenticateClientCredentials(loginURL, clientID, clientSecret string) (instanceURL, accessToken string, err error) {
	if loginURL == "" {
		loginURL = "https://login.salesforce.com"
	}
	tokenURL := strings.TrimRight(loginURL, "/") + "/services/oauth2/token"
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	httpClient := &http.Client{Timeout: 30 * time.Second}
	resp, err := httpClient.PostForm(tokenURL, data)
	if err != nil {
		return "", "", fmt.Errorf("salesforce oauth: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		var authErr types.AuthError
		_ = json.Unmarshal(body, &authErr)
		if authErr.ErrorType != "" {
			authErr.StatusCode = resp.StatusCode
			return "", "", &authErr
		}
		return "", "", fmt.Errorf("salesforce oauth: status %d: %s", resp.StatusCode, body)
	}
	var result struct {
		AccessToken string `json:"access_token"`
		InstanceURL string `json:"instance_url"`
		Error       string `json:"error"`
		Description string `json:"error_description"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", "", fmt.Errorf("salesforce oauth: decode: %w", err)
	}
	if result.Error != "" {
		return "", "", fmt.Errorf("salesforce oauth: %s: %s", result.Error, result.Description)
	}
	return result.InstanceURL, result.AccessToken, nil
}

func tokenEndpoint(loginURL string) string {
	if loginURL == "" {
		return ""
	}
	return strings.TrimRight(loginURL, "/") + "/services/oauth2/token"
}
