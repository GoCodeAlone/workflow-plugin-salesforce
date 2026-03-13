package internal

import (
	"context"
	"fmt"
	"net/url"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// bulkJobPath returns the Bulk API v2 job path.
func bulkJobPath(operation string) string {
	switch operation {
	case "query":
		return "/jobs/query"
	default:
		return "/jobs/ingest"
	}
}

// bulkInsertStep implements step.salesforce_bulk_insert
type bulkInsertStep struct {
	name       string
	moduleName string
}

func newBulkInsertStep(name string, config map[string]any) (*bulkInsertStep, error) {
	return &bulkInsertStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *bulkInsertStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	sObjectType := resolveValue("sobject_type", current, config)
	if sObjectType == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "sobject_type is required"}}, nil
	}
	body := map[string]any{
		"object":    sObjectType,
		"operation": "insert",
	}
	result, err := client.post("/jobs/ingest", body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// bulkUpdateStep implements step.salesforce_bulk_update
type bulkUpdateStep struct {
	name       string
	moduleName string
}

func newBulkUpdateStep(name string, config map[string]any) (*bulkUpdateStep, error) {
	return &bulkUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *bulkUpdateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	sObjectType := resolveValue("sobject_type", current, config)
	if sObjectType == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "sobject_type is required"}}, nil
	}
	body := map[string]any{
		"object":    sObjectType,
		"operation": "update",
	}
	result, err := client.post("/jobs/ingest", body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// bulkUpsertStep implements step.salesforce_bulk_upsert
type bulkUpsertStep struct {
	name       string
	moduleName string
}

func newBulkUpsertStep(name string, config map[string]any) (*bulkUpsertStep, error) {
	return &bulkUpsertStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *bulkUpsertStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	sObjectType := resolveValue("sobject_type", current, config)
	externalIDField := resolveValue("external_id_field", current, config)
	if sObjectType == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "sobject_type is required"}}, nil
	}
	body := map[string]any{
		"object":              sObjectType,
		"operation":           "upsert",
		"externalIdFieldName": externalIDField,
	}
	result, err := client.post("/jobs/ingest", body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// bulkDeleteStep implements step.salesforce_bulk_delete
type bulkDeleteStep struct {
	name       string
	moduleName string
}

func newBulkDeleteStep(name string, config map[string]any) (*bulkDeleteStep, error) {
	return &bulkDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *bulkDeleteStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	sObjectType := resolveValue("sobject_type", current, config)
	if sObjectType == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "sobject_type is required"}}, nil
	}
	body := map[string]any{
		"object":    sObjectType,
		"operation": "delete",
	}
	result, err := client.post("/jobs/ingest", body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// bulkQueryStep implements step.salesforce_bulk_query
type bulkQueryStep struct {
	name       string
	moduleName string
}

func newBulkQueryStep(name string, config map[string]any) (*bulkQueryStep, error) {
	return &bulkQueryStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *bulkQueryStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	soql := resolveValue("soql", current, config)
	if soql == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "soql is required"}}, nil
	}
	body := map[string]any{
		"operation": "query",
		"query":     soql,
	}
	result, err := client.post("/jobs/query", body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// bulkQueryResultsStep implements step.salesforce_bulk_query_results
type bulkQueryResultsStep struct {
	name       string
	moduleName string
}

func newBulkQueryResultsStep(name string, config map[string]any) (*bulkQueryResultsStep, error) {
	return &bulkQueryResultsStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *bulkQueryResultsStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	jobID := resolveValue("job_id", current, config)
	if jobID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "job_id is required"}}, nil
	}
	path := fmt.Sprintf("/jobs/query/%s/results", jobID)
	if maxRecords := resolveInt("max_records", current, config); maxRecords > 0 {
		path += fmt.Sprintf("?maxRecords=%d", maxRecords)
	}
	result, err := client.get(path)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// bulkJobStatusStep implements step.salesforce_bulk_job_status
type bulkJobStatusStep struct {
	name       string
	moduleName string
}

func newBulkJobStatusStep(name string, config map[string]any) (*bulkJobStatusStep, error) {
	return &bulkJobStatusStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *bulkJobStatusStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	jobID := resolveValue("job_id", current, config)
	jobType := resolveValue("job_type", current, config) // "ingest" or "query"
	if jobID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "job_id is required"}}, nil
	}
	basePath := "/jobs/ingest"
	if jobType == "query" {
		basePath = "/jobs/query"
	}
	path := fmt.Sprintf("%s/%s", basePath, jobID)
	result, err := client.get(path)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// bulkJobAbortStep implements step.salesforce_bulk_job_abort
type bulkJobAbortStep struct {
	name       string
	moduleName string
}

func newBulkJobAbortStep(name string, config map[string]any) (*bulkJobAbortStep, error) {
	return &bulkJobAbortStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *bulkJobAbortStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	jobID := resolveValue("job_id", current, config)
	jobType := resolveValue("job_type", current, config)
	if jobID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "job_id is required"}}, nil
	}
	basePath := "/jobs/ingest"
	if jobType == "query" {
		basePath = "/jobs/query"
	}
	path := fmt.Sprintf("%s/%s", basePath, jobID)
	body := map[string]any{"state": "Aborted"}
	result, err := client.patch(path, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// suppress unused
var _ = url.QueryEscape
