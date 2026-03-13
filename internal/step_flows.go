package internal

import (
	"context"
	"fmt"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// flowListStep implements step.salesforce_flow_list
type flowListStep struct {
	name       string
	moduleName string
}

func newFlowListStep(name string, config map[string]any) (*flowListStep, error) {
	return &flowListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *flowListStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	result, err := client.get("/actions/custom/flow")
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// flowRunStep implements step.salesforce_flow_run
type flowRunStep struct {
	name       string
	moduleName string
}

func newFlowRunStep(name string, config map[string]any) (*flowRunStep, error) {
	return &flowRunStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *flowRunStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	flowName := resolveValue("flow_name", current, config)
	if flowName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "flow_name is required"}}, nil
	}
	inputs := resolveMap("inputs", current, config)
	body := map[string]any{}
	if inputs != nil {
		// Salesforce flow invocation expects inputs as list of name/value pairs
		inputList := make([]any, 0, len(inputs))
		for k, v := range inputs {
			inputList = append(inputList, map[string]any{"name": k, "value": v})
		}
		body["inputs"] = []any{map[string]any{"inputs": inputList}}
	}
	path := fmt.Sprintf("/actions/custom/flow/%s", flowName)
	result, err := client.post(path, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}
