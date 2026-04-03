package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	sf "github.com/PramithaMJ/salesforce/v2"
	sfhttp "github.com/PramithaMJ/salesforce/v2/http"
	"github.com/PramithaMJ/salesforce/v2/types"
)

const defaultAPIVersion = "v63.0"

// salesforceClient wraps the PramithaMJ Salesforce SDK HTTP transport
// while preserving the same method signatures all 73+ steps depend on.
type salesforceClient struct {
	sdkHTTP     *sfhttp.Client // SDK HTTP client for authenticated REST calls
	sfClient    *sf.Client     // full SDK client for typed service access (may be nil in tests)
	instanceURL string         // base instance URL — accessed directly by step_apex and step_users
	accessToken string         // current access token
	apiVersion  string         // e.g. "v63.0"
}

// newSalesforceClient creates a client using a pre-existing access token.
// Used for direct-token initialization and test helpers.
func newSalesforceClient(instanceURL, accessToken, apiVersion string) *salesforceClient {
	if apiVersion == "" {
		apiVersion = defaultAPIVersion
	}
	instanceURL = strings.TrimRight(instanceURL, "/")
	hc := sfhttp.NewClient(sfhttp.Config{
		APIVersion: strings.TrimPrefix(apiVersion, "v"),
		MaxRetries: 0,
	})
	hc.SetBaseURL(instanceURL)
	hc.SetAccessToken(accessToken)
	return &salesforceClient{
		sdkHTTP:     hc,
		instanceURL: instanceURL,
		accessToken: accessToken,
		apiVersion:  apiVersion,
	}
}

// newSalesforceClientFromSDK creates a compatibility wrapper around a
// fully-initialised SDK client, keeping the same REST methods that
// step implementations call.
func newSalesforceClientFromSDK(client *sf.Client, apiVersion string) *salesforceClient {
	if apiVersion == "" {
		apiVersion = defaultAPIVersion
	}
	instanceURL := client.InstanceURL()
	token := ""
	if t := client.GetToken(); t != nil {
		token = t.AccessToken
		// WithAccessToken path skips Connect(), so InstanceURL() is empty.
		// Fall back to the token's InstanceURL which was set during auth.
		if instanceURL == "" {
			instanceURL = strings.TrimRight(t.InstanceURL, "/")
		}
	}

	hc := sfhttp.NewClient(sfhttp.Config{
		APIVersion: strings.TrimPrefix(apiVersion, "v"),
		MaxRetries: types.DefaultMaxRetries,
	})
	hc.SetBaseURL(instanceURL)
	hc.SetAccessToken(token)
	return &salesforceClient{
		sdkHTTP:     hc,
		sfClient:    client,
		instanceURL: instanceURL,
		accessToken: token,
		apiVersion:  apiVersion,
	}
}

// versionedPath prepends the versioned REST base path to a relative
// path.  Absolute URLs (starting with "http") are returned as-is.
func (c *salesforceClient) versionedPath(path string) string {
	if strings.HasPrefix(path, "http") {
		// Absolute URL — strip the instance URL prefix so the SDK HTTP
		// client (which prepends its own baseURL) produces the correct
		// final URL.
		if strings.HasPrefix(path, c.instanceURL) {
			return strings.TrimPrefix(path, c.instanceURL)
		}
		return path
	}
	return fmt.Sprintf("/services/data/%s%s", c.apiVersion, path)
}

// baseURL returns the versioned API base URL (kept for backward compat).
func (c *salesforceClient) baseURL() string {
	return fmt.Sprintf("%s/services/data/%s", c.instanceURL, c.apiVersion)
}

// do performs an HTTP request via the SDK transport and decodes the
// JSON response.  The method signature is unchanged from the original
// custom client so all step implementations continue to compile.
func (c *salesforceClient) do(method, path string, body any) (map[string]any, int, error) {
	sdkPath := c.versionedPath(path)
	ctx := context.Background()

	var respBody []byte
	var err error

	switch method {
	case http.MethodGet:
		respBody, err = c.sdkHTTP.Get(ctx, sdkPath)
	case http.MethodPost:
		respBody, err = c.sdkHTTP.Post(ctx, sdkPath, body)
	case http.MethodPatch:
		respBody, err = c.sdkHTTP.Patch(ctx, sdkPath, body)
	case http.MethodDelete:
		respBody, err = c.sdkHTTP.Delete(ctx, sdkPath)
	case http.MethodPut:
		respBody, err = c.sdkHTTP.Put(ctx, sdkPath, body)
	default:
		return nil, 0, fmt.Errorf("salesforce: unsupported method: %s", method)
	}

	if err != nil {
		statusCode := extractStatusCode(err)
		return nil, statusCode, fmt.Errorf("salesforce: %s", err.Error())
	}

	// Empty response (e.g. 204 No Content)
	if len(respBody) == 0 {
		return map[string]any{"success": true}, 204, nil
	}

	// Salesforce sometimes returns an array
	if respBody[0] == '[' {
		var arr []map[string]any
		if err := json.Unmarshal(respBody, &arr); err != nil {
			return nil, 200, fmt.Errorf("salesforce: decode array: %w", err)
		}
		if len(arr) > 0 {
			return arr[0], 200, nil
		}
		return map[string]any{}, 200, nil
	}

	var result map[string]any
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, 200, fmt.Errorf("salesforce: decode response: %w", err)
	}
	return result, 200, nil
}

// doArray performs an HTTP request and decodes the response as an
// array (for collection endpoints, SOQL results, etc.).
func (c *salesforceClient) doArray(method, path string, body any) ([]any, map[string]any, int, error) {
	sdkPath := c.versionedPath(path)
	ctx := context.Background()

	var respBody []byte
	var err error

	switch method {
	case http.MethodGet:
		respBody, err = c.sdkHTTP.Get(ctx, sdkPath)
	case http.MethodPost:
		respBody, err = c.sdkHTTP.Post(ctx, sdkPath, body)
	case http.MethodPatch:
		respBody, err = c.sdkHTTP.Patch(ctx, sdkPath, body)
	case http.MethodDelete:
		respBody, err = c.sdkHTTP.Delete(ctx, sdkPath)
	default:
		return nil, nil, 0, fmt.Errorf("salesforce: unsupported method: %s", method)
	}

	if err != nil {
		statusCode := extractStatusCode(err)
		return nil, nil, statusCode, fmt.Errorf("salesforce: %s", err.Error())
	}

	if len(respBody) == 0 {
		return nil, map[string]any{}, 200, nil
	}

	// Try array first
	if respBody[0] == '[' {
		var arr []any
		if err := json.Unmarshal(respBody, &arr); err != nil {
			return nil, nil, 200, fmt.Errorf("salesforce: decode array: %w", err)
		}
		return arr, nil, 200, nil
	}

	// Fall back to object
	var result map[string]any
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, nil, 200, fmt.Errorf("salesforce: decode response: %w", err)
	}
	return nil, result, 200, nil
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

// extractStatusCode pulls an HTTP status code from SDK error types.
func extractStatusCode(err error) int {
	if apiErr, ok := err.(*types.APIError); ok {
		return apiErr.StatusCode
	}
	if apiErrs, ok := err.(types.APIErrors); ok && len(apiErrs) > 0 {
		return apiErrs[0].StatusCode
	}
	if authErr, ok := err.(*types.AuthError); ok {
		return authErr.StatusCode
	}
	return 500
}
