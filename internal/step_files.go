package internal

import (
	"context"
	"encoding/base64"
	"fmt"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// fileUploadStep implements step.salesforce_file_upload
type fileUploadStep struct {
	name       string
	moduleName string
}

func newFileUploadStep(name string, config map[string]any) (*fileUploadStep, error) {
	return &fileUploadStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *fileUploadStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	title := resolveValue("title", current, config)
	if title == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "title is required"}}, nil
	}
	// file_content should be base64-encoded or raw bytes as string
	fileContent := resolveValue("file_content", current, config)
	pathOnClient := resolveValue("path_on_client", current, config)
	versionData := base64.StdEncoding.EncodeToString([]byte(fileContent))
	body := map[string]any{
		"Title":        title,
		"PathOnClient": pathOnClient,
		"VersionData":  versionData,
	}
	result, err := client.post("/sobjects/ContentVersion", body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// fileDownloadStep implements step.salesforce_file_download
type fileDownloadStep struct {
	name       string
	moduleName string
}

func newFileDownloadStep(name string, config map[string]any) (*fileDownloadStep, error) {
	return &fileDownloadStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *fileDownloadStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	contentVersionID := resolveValue("content_version_id", current, config)
	if contentVersionID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "content_version_id is required"}}, nil
	}
	path := fmt.Sprintf("/sobjects/ContentVersion/%s/VersionData", contentVersionID)
	result, err := client.get(path)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// contentVersionCreateStep implements step.salesforce_content_version_create
type contentVersionCreateStep struct {
	name       string
	moduleName string
}

func newContentVersionCreateStep(name string, config map[string]any) (*contentVersionCreateStep, error) {
	return &contentVersionCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *contentVersionCreateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	fields := resolveMap("fields", current, config)
	if fields == nil {
		fields = map[string]any{}
	}
	result, err := client.post("/sobjects/ContentVersion", fields)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// contentDocumentGetStep implements step.salesforce_content_document_get
type contentDocumentGetStep struct {
	name       string
	moduleName string
}

func newContentDocumentGetStep(name string, config map[string]any) (*contentDocumentGetStep, error) {
	return &contentDocumentGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *contentDocumentGetStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	documentID := resolveValue("document_id", current, config)
	if documentID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "document_id is required"}}, nil
	}
	path := fmt.Sprintf("/sobjects/ContentDocument/%s", documentID)
	result, err := client.get(path)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// contentDocumentDeleteStep implements step.salesforce_content_document_delete
type contentDocumentDeleteStep struct {
	name       string
	moduleName string
}

func newContentDocumentDeleteStep(name string, config map[string]any) (*contentDocumentDeleteStep, error) {
	return &contentDocumentDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *contentDocumentDeleteStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	documentID := resolveValue("document_id", current, config)
	if documentID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "document_id is required"}}, nil
	}
	path := fmt.Sprintf("/sobjects/ContentDocument/%s", documentID)
	result, err := client.delete(path)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}
