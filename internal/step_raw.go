package internal

import (
	"context"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// rawRequestStep implements step.salesforce_raw_request
type rawRequestStep struct {
	name       string
	moduleName string
}

func newRawRequestStep(name string, config map[string]any) (*rawRequestStep, error) {
	return &rawRequestStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *rawRequestStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	method := resolveValue("method", current, config)
	path := resolveValue("path", current, config)
	if method == "" {
		method = "GET"
	}
	if path == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "path is required"}}, nil
	}
	body := resolveMap("body", current, config)
	result, _, err := client.do(method, path, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}
