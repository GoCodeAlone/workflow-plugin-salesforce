package internal

import (
	"context"
	"fmt"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// metadataDescribeStep implements step.salesforce_metadata_describe
type metadataDescribeStep struct {
	name       string
	moduleName string
}

func newMetadataDescribeStep(name string, config map[string]any) (*metadataDescribeStep, error) {
	return &metadataDescribeStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *metadataDescribeStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	result, err := client.get("/tooling/sobjects")
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// metadataListStep implements step.salesforce_metadata_list
type metadataListStep struct {
	name       string
	moduleName string
}

func newMetadataListStep(name string, config map[string]any) (*metadataListStep, error) {
	return &metadataListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *metadataListStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	metadataType := resolveValue("metadata_type", current, config)
	if metadataType == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "metadata_type is required"}}, nil
	}
	path := fmt.Sprintf("/tooling/query?q=SELECT+Id,FullName+FROM+%s", metadataType)
	result, err := client.get(path)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// metadataReadStep implements step.salesforce_metadata_read
type metadataReadStep struct {
	name       string
	moduleName string
}

func newMetadataReadStep(name string, config map[string]any) (*metadataReadStep, error) {
	return &metadataReadStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *metadataReadStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	metadataType := resolveValue("metadata_type", current, config)
	recordID := resolveValue("record_id", current, config)
	if metadataType == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "metadata_type is required"}}, nil
	}
	if recordID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "record_id is required"}}, nil
	}
	path := fmt.Sprintf("/tooling/sobjects/%s/%s", metadataType, recordID)
	result, err := client.get(path)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// metadataCreateStep implements step.salesforce_metadata_create
type metadataCreateStep struct {
	name       string
	moduleName string
}

func newMetadataCreateStep(name string, config map[string]any) (*metadataCreateStep, error) {
	return &metadataCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *metadataCreateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	metadataType := resolveValue("metadata_type", current, config)
	if metadataType == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "metadata_type is required"}}, nil
	}
	fields := resolveMap("fields", current, config)
	if fields == nil {
		fields = map[string]any{}
	}
	path := fmt.Sprintf("/tooling/sobjects/%s", metadataType)
	result, err := client.post(path, fields)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// metadataUpdateStep implements step.salesforce_metadata_update
type metadataUpdateStep struct {
	name       string
	moduleName string
}

func newMetadataUpdateStep(name string, config map[string]any) (*metadataUpdateStep, error) {
	return &metadataUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *metadataUpdateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	metadataType := resolveValue("metadata_type", current, config)
	recordID := resolveValue("record_id", current, config)
	if metadataType == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "metadata_type is required"}}, nil
	}
	if recordID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "record_id is required"}}, nil
	}
	fields := resolveMap("fields", current, config)
	if fields == nil {
		fields = map[string]any{}
	}
	path := fmt.Sprintf("/tooling/sobjects/%s/%s", metadataType, recordID)
	result, err := client.patch(path, fields)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// metadataDeleteStep implements step.salesforce_metadata_delete
type metadataDeleteStep struct {
	name       string
	moduleName string
}

func newMetadataDeleteStep(name string, config map[string]any) (*metadataDeleteStep, error) {
	return &metadataDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *metadataDeleteStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	metadataType := resolveValue("metadata_type", current, config)
	recordID := resolveValue("record_id", current, config)
	if metadataType == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "metadata_type is required"}}, nil
	}
	if recordID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "record_id is required"}}, nil
	}
	path := fmt.Sprintf("/tooling/sobjects/%s/%s", metadataType, recordID)
	result, err := client.delete(path)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// metadataDeployStep implements step.salesforce_metadata_deploy
type metadataDeployStep struct {
	name       string
	moduleName string
}

func newMetadataDeployStep(name string, config map[string]any) (*metadataDeployStep, error) {
	return &metadataDeployStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *metadataDeployStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	deploymentBody := resolveMap("deployment", current, config)
	if deploymentBody == nil {
		return &sdk.StepResult{Output: map[string]any{"error": "deployment is required"}}, nil
	}
	result, err := client.post("/metadata/deployRequest", deploymentBody)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// metadataRetrieveStep implements step.salesforce_metadata_retrieve
type metadataRetrieveStep struct {
	name       string
	moduleName string
}

func newMetadataRetrieveStep(name string, config map[string]any) (*metadataRetrieveStep, error) {
	return &metadataRetrieveStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *metadataRetrieveStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	retrieveBody := resolveMap("retrieve", current, config)
	if retrieveBody == nil {
		return &sdk.StepResult{Output: map[string]any{"error": "retrieve is required"}}, nil
	}
	result, err := client.post("/metadata/retrieveRequest", retrieveBody)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}
