package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultAPIVersion = "v63.0"

// salesforceClient is a lightweight HTTP REST client for the Salesforce REST API.
type salesforceClient struct {
	httpClient  *http.Client
	instanceURL string
	accessToken string
	apiVersion  string
}

// newSalesforceClient creates a client using a pre-existing access token.
func newSalesforceClient(instanceURL, accessToken, apiVersion string) *salesforceClient {
	if apiVersion == "" {
		apiVersion = defaultAPIVersion
	}
	return &salesforceClient{
		httpClient:  &http.Client{Timeout: 30 * time.Second},
		instanceURL: strings.TrimRight(instanceURL, "/"),
		accessToken: accessToken,
		apiVersion:  apiVersion,
	}
}

// authenticate uses OAuth2 client credentials (or password flow) to obtain an access token.
func authenticateOAuth(loginURL, clientID, clientSecret string) (instanceURL, accessToken string, err error) {
	if loginURL == "" {
		loginURL = "https://login.salesforce.com"
	}
	tokenURL := strings.TrimRight(loginURL, "/") + "/services/oauth2/token"
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		return "", "", fmt.Errorf("salesforce oauth: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
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

// baseURL returns the versioned API base URL.
func (c *salesforceClient) baseURL() string {
	return fmt.Sprintf("%s/services/data/%s", c.instanceURL, c.apiVersion)
}

// do performs an HTTP request and decodes the JSON response.
func (c *salesforceClient) do(method, path string, body any) (map[string]any, int, error) {
	fullURL := path
	if !strings.HasPrefix(path, "http") {
		fullURL = c.baseURL() + path
	}

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("salesforce: marshal request: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, fullURL, reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("salesforce: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("salesforce: http: %w", err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusNoContent {
		return map[string]any{"success": true}, resp.StatusCode, nil
	}

	// Salesforce sometimes returns an array for error responses
	if len(respBody) > 0 && respBody[0] == '[' {
		var arr []map[string]any
		if err := json.Unmarshal(respBody, &arr); err != nil {
			return nil, resp.StatusCode, fmt.Errorf("salesforce: decode array: %w", err)
		}
		if len(arr) > 0 {
			return arr[0], resp.StatusCode, nil
		}
		return map[string]any{}, resp.StatusCode, nil
	}

	var result map[string]any
	if len(respBody) > 0 {
		if err := json.Unmarshal(respBody, &result); err != nil {
			return nil, resp.StatusCode, fmt.Errorf("salesforce: decode response: %w", err)
		}
	}
	if result == nil {
		result = map[string]any{}
	}

	if resp.StatusCode >= 400 {
		errMsg := fmt.Sprintf("status %d", resp.StatusCode)
		if msg, ok := result["message"].(string); ok {
			errMsg = msg
		} else if errCode, ok := result["errorCode"].(string); ok {
			errMsg = errCode
		}
		return result, resp.StatusCode, fmt.Errorf("salesforce: %s", errMsg)
	}

	return result, resp.StatusCode, nil
}

// doArray performs an HTTP GET and decodes the response as an array (for SOQL results etc.).
func (c *salesforceClient) doArray(method, path string, body any) ([]any, map[string]any, int, error) {
	fullURL := path
	if !strings.HasPrefix(path, "http") {
		fullURL = c.baseURL() + path
	}

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("salesforce: marshal request: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, fullURL, reqBody)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("salesforce: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("salesforce: http: %w", err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)

	// Try array first
	if len(respBody) > 0 && respBody[0] == '[' {
		var arr []any
		if err := json.Unmarshal(respBody, &arr); err != nil {
			return nil, nil, resp.StatusCode, fmt.Errorf("salesforce: decode array: %w", err)
		}
		return arr, nil, resp.StatusCode, nil
	}

	// Fall back to object
	var result map[string]any
	if len(respBody) > 0 {
		if err := json.Unmarshal(respBody, &result); err != nil {
			return nil, nil, resp.StatusCode, fmt.Errorf("salesforce: decode response: %w", err)
		}
	}
	if resp.StatusCode >= 400 {
		errMsg := fmt.Sprintf("status %d", resp.StatusCode)
		if result != nil {
			if msg, ok := result["message"].(string); ok {
				errMsg = msg
			}
		}
		return nil, result, resp.StatusCode, fmt.Errorf("salesforce: %s", errMsg)
	}
	return nil, result, resp.StatusCode, nil
}

// get is a convenience wrapper for GET requests.
func (c *salesforceClient) get(path string) (map[string]any, error) {
	result, _, err := c.do(http.MethodGet, path, nil)
	return result, err
}

// post is a convenience wrapper for POST requests.
func (c *salesforceClient) post(path string, body any) (map[string]any, error) {
	result, _, err := c.do(http.MethodPost, path, body)
	return result, err
}

// patch is a convenience wrapper for PATCH requests.
func (c *salesforceClient) patch(path string, body any) (map[string]any, error) {
	result, _, err := c.do(http.MethodPatch, path, body)
	return result, err
}

// delete is a convenience wrapper for DELETE requests.
func (c *salesforceClient) delete(path string) (map[string]any, error) {
	result, _, err := c.do(http.MethodDelete, path, nil)
	return result, err
}

// getArray returns the response parsed as an array.
func (c *salesforceClient) getArray(path string) ([]any, map[string]any, error) {
	arr, obj, _, err := c.doArray(http.MethodGet, path, nil)
	return arr, obj, err
}

// postArray posts and returns array response.
func (c *salesforceClient) postArray(path string, body any) ([]any, map[string]any, error) {
	arr, obj, _, err := c.doArray(http.MethodPost, path, body)
	return arr, obj, err
}
