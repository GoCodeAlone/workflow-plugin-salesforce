package internal

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// newTestClient creates a salesforceClient pointing at the given test server.
func newTestClient(srv *httptest.Server) *salesforceClient {
	return &salesforceClient{
		httpClient:  srv.Client(),
		instanceURL: srv.URL,
		accessToken: "test-token",
		apiVersion:  defaultAPIVersion,
	}
}

func TestRecordGetStep_MissingClient(t *testing.T) {
	step, err := newRecordGetStep("test", map[string]any{"module": "nonexistent"})
	if err != nil {
		t.Fatal(err)
	}
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{"sobject_type": "Account", "record_id": "001xx"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] == nil {
		t.Error("expected error for missing client")
	}
}

func TestRecordGetStep_MissingSObjectType(t *testing.T) {
	step, _ := newRecordGetStep("test", map[string]any{"module": "nonexistent2"})
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] == nil {
		t.Error("expected error for missing sobject_type")
	}
}

func TestRecordGetStep_MissingRecordID(t *testing.T) {
	step, _ := newRecordGetStep("test", map[string]any{"module": "nonexistent3"})
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{"sobject_type": "Account"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] == nil {
		t.Error("expected error for missing record_id")
	}
}

func TestRecordGetStep_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/services/data/v63.0/sobjects/Account/001xx" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"Id": "001xx", "Name": "Acme Corp"})
	}))
	defer srv.Close()

	RegisterClient("test-get-ok", newTestClient(srv))
	defer UnregisterClient("test-get-ok")

	step, _ := newRecordGetStep("test", map[string]any{"module": "test-get-ok"})
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{
		"sobject_type": "Account",
		"record_id":    "001xx",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] != nil {
		t.Fatalf("unexpected error: %v", result.Output["error"])
	}
	if result.Output["Id"] != "001xx" {
		t.Errorf("expected Id=001xx, got %v", result.Output["Id"])
	}
}

func TestRecordCreateStep_MissingClient(t *testing.T) {
	step, _ := newRecordCreateStep("test", map[string]any{"module": "nonexistent-create"})
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{"sobject_type": "Account"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] == nil {
		t.Error("expected error for missing client")
	}
}

func TestRecordCreateStep_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{"id": "001newid", "success": true})
	}))
	defer srv.Close()

	RegisterClient("test-create-ok", newTestClient(srv))
	defer UnregisterClient("test-create-ok")

	step, _ := newRecordCreateStep("test", map[string]any{"module": "test-create-ok"})
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{
		"sobject_type": "Account",
		"fields":       map[string]any{"Name": "Test Account"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] != nil {
		t.Fatalf("unexpected error: %v", result.Output["error"])
	}
}

func TestRecordDeleteStep_MissingClient(t *testing.T) {
	step, _ := newRecordDeleteStep("test", map[string]any{"module": "nonexistent-del"})
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{
		"sobject_type": "Account",
		"record_id":    "001xx",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] == nil {
		t.Error("expected error for missing client")
	}
}

func TestDescribeGlobalStep_MissingClient(t *testing.T) {
	step, _ := newDescribeGlobalStep("test", map[string]any{"module": "nonexistent-describe"})
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] == nil {
		t.Error("expected error for missing client")
	}
}

func TestQueryStep_MissingSOQL(t *testing.T) {
	step, _ := newQueryStep("test", map[string]any{"module": "nonexistent-query"})
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] == nil {
		t.Error("expected error for missing soql")
	}
}

func TestQueryStep_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"totalSize": 2,
			"done":      true,
			"records":   []any{map[string]any{"Id": "001a"}, map[string]any{"Id": "001b"}},
		})
	}))
	defer srv.Close()

	RegisterClient("test-query-ok", newTestClient(srv))
	defer UnregisterClient("test-query-ok")

	step, _ := newQueryStep("test", map[string]any{"module": "test-query-ok"})
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{
		"soql": "SELECT Id FROM Account",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] != nil {
		t.Fatalf("unexpected error: %v", result.Output["error"])
	}
	if result.Output["total_size"].(int) != 2 {
		t.Errorf("expected total_size=2, got %v", result.Output["total_size"])
	}
}

func TestRawRequestStep_MissingPath(t *testing.T) {
	step, _ := newRawRequestStep("test", map[string]any{"module": "nonexistent-raw"})
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] == nil {
		t.Error("expected error for missing path")
	}
}

func TestStepRegistry_AllTypesConstructible(t *testing.T) {
	for typeName := range stepRegistry {
		_, err := createStep(typeName, "test-"+typeName, map[string]any{})
		if err != nil {
			t.Errorf("createStep(%q): unexpected error: %v", typeName, err)
		}
	}
}

func TestStepRegistry_UnknownType(t *testing.T) {
	_, err := createStep("step.salesforce_unknown_xyz", "test", map[string]any{})
	if err == nil {
		t.Error("expected error for unknown step type")
	}
}
