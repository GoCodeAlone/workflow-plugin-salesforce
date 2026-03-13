package internal

import (
	"context"
	"net/url"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// queryStep implements step.salesforce_query (SOQL)
type queryStep struct {
	name       string
	moduleName string
}

func newQueryStep(name string, config map[string]any) (*queryStep, error) {
	return &queryStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *queryStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	soql := resolveValue("soql", current, config)
	if soql == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "soql is required"}}, nil
	}
	path := "/query?q=" + url.QueryEscape(soql)
	result, err := client.get(path)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	records, _ := result["records"].([]any)
	totalSize, _ := result["totalSize"].(float64)
	done, _ := result["done"].(bool)
	return &sdk.StepResult{Output: map[string]any{
		"records":    records,
		"total_size": int(totalSize),
		"done":       done,
		"next_url":   result["nextRecordsUrl"],
	}}, nil
}

// queryAllStep implements step.salesforce_query_all (includes deleted/archived records)
type queryAllStep struct {
	name       string
	moduleName string
}

func newQueryAllStep(name string, config map[string]any) (*queryAllStep, error) {
	return &queryAllStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *queryAllStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	soql := resolveValue("soql", current, config)
	if soql == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "soql is required"}}, nil
	}
	path := "/queryAll?q=" + url.QueryEscape(soql)
	result, err := client.get(path)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	records, _ := result["records"].([]any)
	totalSize, _ := result["totalSize"].(float64)
	done, _ := result["done"].(bool)
	return &sdk.StepResult{Output: map[string]any{
		"records":    records,
		"total_size": int(totalSize),
		"done":       done,
		"next_url":   result["nextRecordsUrl"],
	}}, nil
}

// searchStep implements step.salesforce_search (SOSL)
type searchStep struct {
	name       string
	moduleName string
}

func newSearchStep(name string, config map[string]any) (*searchStep, error) {
	return &searchStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *searchStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	sosl := resolveValue("sosl", current, config)
	if sosl == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "sosl is required"}}, nil
	}
	path := "/search?q=" + url.QueryEscape(sosl)
	result, err := client.get(path)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}
