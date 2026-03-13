package internal

import (
	"context"
	"fmt"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// recordGetStep implements step.salesforce_record_get
type recordGetStep struct {
	name       string
	moduleName string
}

func newRecordGetStep(name string, config map[string]any) (*recordGetStep, error) {
	return &recordGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *recordGetStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
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
	path := fmt.Sprintf("/sobjects/%s/%s", sObjectType, recordID)
	result, err := client.get(path)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// recordCreateStep implements step.salesforce_record_create
type recordCreateStep struct {
	name       string
	moduleName string
}

func newRecordCreateStep(name string, config map[string]any) (*recordCreateStep, error) {
	return &recordCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *recordCreateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
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
	path := fmt.Sprintf("/sobjects/%s", sObjectType)
	result, err := client.post(path, fields)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// recordUpdateStep implements step.salesforce_record_update
type recordUpdateStep struct {
	name       string
	moduleName string
}

func newRecordUpdateStep(name string, config map[string]any) (*recordUpdateStep, error) {
	return &recordUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *recordUpdateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
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
	path := fmt.Sprintf("/sobjects/%s/%s", sObjectType, recordID)
	result, err := client.patch(path, fields)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// recordUpsertStep implements step.salesforce_record_upsert
type recordUpsertStep struct {
	name       string
	moduleName string
}

func newRecordUpsertStep(name string, config map[string]any) (*recordUpsertStep, error) {
	return &recordUpsertStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *recordUpsertStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	sObjectType := resolveValue("sobject_type", current, config)
	externalIDField := resolveValue("external_id_field", current, config)
	externalIDValue := resolveValue("external_id_value", current, config)
	if sObjectType == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "sobject_type is required"}}, nil
	}
	if externalIDField == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "external_id_field is required"}}, nil
	}
	if externalIDValue == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "external_id_value is required"}}, nil
	}
	fields := resolveMap("fields", current, config)
	if fields == nil {
		fields = map[string]any{}
	}
	path := fmt.Sprintf("/sobjects/%s/%s/%s", sObjectType, externalIDField, externalIDValue)
	result, err := client.patch(path, fields)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// recordDeleteStep implements step.salesforce_record_delete
type recordDeleteStep struct {
	name       string
	moduleName string
}

func newRecordDeleteStep(name string, config map[string]any) (*recordDeleteStep, error) {
	return &recordDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *recordDeleteStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
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
	path := fmt.Sprintf("/sobjects/%s/%s", sObjectType, recordID)
	result, err := client.delete(path)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// recordDescribeStep implements step.salesforce_record_describe
type recordDescribeStep struct {
	name       string
	moduleName string
}

func newRecordDescribeStep(name string, config map[string]any) (*recordDescribeStep, error) {
	return &recordDescribeStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *recordDescribeStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	sObjectType := resolveValue("sobject_type", current, config)
	if sObjectType == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "sobject_type is required"}}, nil
	}
	path := fmt.Sprintf("/sobjects/%s/describe", sObjectType)
	result, err := client.get(path)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// describeGlobalStep implements step.salesforce_describe_global
type describeGlobalStep struct {
	name       string
	moduleName string
}

func newDescribeGlobalStep(name string, config map[string]any) (*describeGlobalStep, error) {
	return &describeGlobalStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *describeGlobalStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	result, err := client.get("/sobjects")
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}
