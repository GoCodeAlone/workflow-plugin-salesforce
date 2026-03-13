package internal

import (
	"context"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// compositeRequestStep implements step.salesforce_composite_request
type compositeRequestStep struct {
	name       string
	moduleName string
}

func newCompositeRequestStep(name string, config map[string]any) (*compositeRequestStep, error) {
	return &compositeRequestStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *compositeRequestStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	compositeReqs := resolveAnySlice("composite_request", current, config)
	if len(compositeReqs) == 0 {
		return &sdk.StepResult{Output: map[string]any{"error": "composite_request is required"}}, nil
	}
	allOrNone := resolveBool("all_or_none", current, config)
	body := map[string]any{
		"compositeRequest": compositeReqs,
		"allOrNone":        allOrNone,
	}
	result, err := client.post("/composite", body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// compositeTreeStep implements step.salesforce_composite_tree
type compositeTreeStep struct {
	name       string
	moduleName string
}

func newCompositeTreeStep(name string, config map[string]any) (*compositeTreeStep, error) {
	return &compositeTreeStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *compositeTreeStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	sObjectType := resolveValue("sobject_type", current, config)
	records := resolveAnySlice("records", current, config)
	if sObjectType == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "sobject_type is required"}}, nil
	}
	if len(records) == 0 {
		return &sdk.StepResult{Output: map[string]any{"error": "records is required"}}, nil
	}
	body := map[string]any{"records": records}
	path := "/composite/tree/" + sObjectType
	result, err := client.post(path, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}
