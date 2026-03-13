package internal

import (
	"context"
	"fmt"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// apexGetStep implements step.salesforce_apex_get
type apexGetStep struct {
	name       string
	moduleName string
}

func newApexGetStep(name string, config map[string]any) (*apexGetStep, error) {
	return &apexGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *apexGetStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	apexPath := resolveValue("apex_path", current, config)
	if apexPath == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "apex_path is required"}}, nil
	}
	path := fmt.Sprintf("%s/services/apexrest%s", client.instanceURL, apexPath)
	result, err := client.get(path)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// apexPostStep implements step.salesforce_apex_post
type apexPostStep struct {
	name       string
	moduleName string
}

func newApexPostStep(name string, config map[string]any) (*apexPostStep, error) {
	return &apexPostStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *apexPostStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	apexPath := resolveValue("apex_path", current, config)
	if apexPath == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "apex_path is required"}}, nil
	}
	body := resolveMap("body", current, config)
	path := fmt.Sprintf("%s/services/apexrest%s", client.instanceURL, apexPath)
	result, err := client.post(path, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// apexPatchStep implements step.salesforce_apex_patch
type apexPatchStep struct {
	name       string
	moduleName string
}

func newApexPatchStep(name string, config map[string]any) (*apexPatchStep, error) {
	return &apexPatchStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *apexPatchStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	apexPath := resolveValue("apex_path", current, config)
	if apexPath == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "apex_path is required"}}, nil
	}
	body := resolveMap("body", current, config)
	path := fmt.Sprintf("%s/services/apexrest%s", client.instanceURL, apexPath)
	result, err := client.patch(path, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// apexPutStep implements step.salesforce_apex_put
type apexPutStep struct {
	name       string
	moduleName string
}

func newApexPutStep(name string, config map[string]any) (*apexPutStep, error) {
	return &apexPutStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *apexPutStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	apexPath := resolveValue("apex_path", current, config)
	if apexPath == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "apex_path is required"}}, nil
	}
	body := resolveMap("body", current, config)
	path := fmt.Sprintf("%s/services/apexrest%s", client.instanceURL, apexPath)
	result, _, err := client.do("PUT", path, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// apexDeleteStep implements step.salesforce_apex_delete
type apexDeleteStep struct {
	name       string
	moduleName string
}

func newApexDeleteStep(name string, config map[string]any) (*apexDeleteStep, error) {
	return &apexDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *apexDeleteStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	apexPath := resolveValue("apex_path", current, config)
	if apexPath == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "apex_path is required"}}, nil
	}
	path := fmt.Sprintf("%s/services/apexrest%s", client.instanceURL, apexPath)
	result, err := client.delete(path)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}
