package internal

import (
	"context"
	"fmt"
	"net/url"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// toolingQueryStep implements step.salesforce_tooling_query
type toolingQueryStep struct {
	name       string
	moduleName string
}

func newToolingQueryStep(name string, config map[string]any) (*toolingQueryStep, error) {
	return &toolingQueryStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *toolingQueryStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	soql := resolveValue("soql", current, config)
	if soql == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "soql is required"}}, nil
	}
	path := "/tooling/query?q=" + url.QueryEscape(soql)
	result, err := client.get(path)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// toolingGetStep implements step.salesforce_tooling_get
type toolingGetStep struct {
	name       string
	moduleName string
}

func newToolingGetStep(name string, config map[string]any) (*toolingGetStep, error) {
	return &toolingGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *toolingGetStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	sObjectType := resolveValue("sobject_type", current, config)
	recordID := resolveValue("record_id", current, config)
	if sObjectType == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "sobject_type is required"}}, nil
	}
	if recordID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "record_id is required"}}, nil
	}
	path := fmt.Sprintf("/tooling/sobjects/%s/%s", sObjectType, recordID)
	result, err := client.get(path)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// toolingCreateStep implements step.salesforce_tooling_create
type toolingCreateStep struct {
	name       string
	moduleName string
}

func newToolingCreateStep(name string, config map[string]any) (*toolingCreateStep, error) {
	return &toolingCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *toolingCreateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	sObjectType := resolveValue("sobject_type", current, config)
	if sObjectType == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "sobject_type is required"}}, nil
	}
	fields := resolveMap("fields", current, config)
	if fields == nil {
		fields = map[string]any{}
	}
	path := fmt.Sprintf("/tooling/sobjects/%s", sObjectType)
	result, err := client.post(path, fields)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// toolingUpdateStep implements step.salesforce_tooling_update
type toolingUpdateStep struct {
	name       string
	moduleName string
}

func newToolingUpdateStep(name string, config map[string]any) (*toolingUpdateStep, error) {
	return &toolingUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *toolingUpdateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	sObjectType := resolveValue("sobject_type", current, config)
	recordID := resolveValue("record_id", current, config)
	if sObjectType == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "sobject_type is required"}}, nil
	}
	if recordID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "record_id is required"}}, nil
	}
	fields := resolveMap("fields", current, config)
	if fields == nil {
		fields = map[string]any{}
	}
	path := fmt.Sprintf("/tooling/sobjects/%s/%s", sObjectType, recordID)
	result, err := client.patch(path, fields)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// toolingDeleteStep implements step.salesforce_tooling_delete
type toolingDeleteStep struct {
	name       string
	moduleName string
}

func newToolingDeleteStep(name string, config map[string]any) (*toolingDeleteStep, error) {
	return &toolingDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *toolingDeleteStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	sObjectType := resolveValue("sobject_type", current, config)
	recordID := resolveValue("record_id", current, config)
	if sObjectType == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "sobject_type is required"}}, nil
	}
	if recordID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "record_id is required"}}, nil
	}
	path := fmt.Sprintf("/tooling/sobjects/%s/%s", sObjectType, recordID)
	result, err := client.delete(path)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// apexExecuteStep implements step.salesforce_apex_execute (anonymous Apex via Tooling API)
type apexExecuteStep struct {
	name       string
	moduleName string
}

func newApexExecuteStep(name string, config map[string]any) (*apexExecuteStep, error) {
	return &apexExecuteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *apexExecuteStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	apexCode := resolveValue("apex_code", current, config)
	if apexCode == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "apex_code is required"}}, nil
	}
	path := "/tooling/executeAnonymous?anonymousBody=" + url.QueryEscape(apexCode)
	result, err := client.get(path)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}
