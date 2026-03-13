package internal

import (
	"context"
	"fmt"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// eventPublishStep implements step.salesforce_event_publish (Platform Events)
type eventPublishStep struct {
	name       string
	moduleName string
}

func newEventPublishStep(name string, config map[string]any) (*eventPublishStep, error) {
	return &eventPublishStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *eventPublishStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	eventType := resolveValue("event_type", current, config)
	if eventType == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "event_type is required"}}, nil
	}
	fields := resolveMap("fields", current, config)
	if fields == nil {
		fields = map[string]any{}
	}
	path := fmt.Sprintf("/sobjects/%s", eventType)
	result, err := client.post(path, fields)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}
