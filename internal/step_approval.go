package internal

import (
	"context"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// approvalListStep implements step.salesforce_approval_list
type approvalListStep struct {
	name       string
	moduleName string
}

func newApprovalListStep(name string, config map[string]any) (*approvalListStep, error) {
	return &approvalListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *approvalListStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	result, err := client.get("/process/approvals")
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// approvalSubmitStep implements step.salesforce_approval_submit
type approvalSubmitStep struct {
	name       string
	moduleName string
}

func newApprovalSubmitStep(name string, config map[string]any) (*approvalSubmitStep, error) {
	return &approvalSubmitStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *approvalSubmitStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	recordID := resolveValue("record_id", current, config)
	if recordID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "record_id is required"}}, nil
	}
	comments := resolveValue("comments", current, config)
	body := map[string]any{
		"requests": []any{
			map[string]any{
				"actionType": "Submit",
				"contextId":  recordID,
				"comments":   comments,
			},
		},
	}
	result, err := client.post("/process/approvals", body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// approvalApproveStep implements step.salesforce_approval_approve
type approvalApproveStep struct {
	name       string
	moduleName string
}

func newApprovalApproveStep(name string, config map[string]any) (*approvalApproveStep, error) {
	return &approvalApproveStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *approvalApproveStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	workItemID := resolveValue("work_item_id", current, config)
	if workItemID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "work_item_id is required"}}, nil
	}
	comments := resolveValue("comments", current, config)
	body := map[string]any{
		"requests": []any{
			map[string]any{
				"actionType": "Approve",
				"workitemId": workItemID,
				"comments":   comments,
			},
		},
	}
	result, err := client.post("/process/approvals", body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// approvalRejectStep implements step.salesforce_approval_reject
type approvalRejectStep struct {
	name       string
	moduleName string
}

func newApprovalRejectStep(name string, config map[string]any) (*approvalRejectStep, error) {
	return &approvalRejectStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *approvalRejectStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	workItemID := resolveValue("work_item_id", current, config)
	if workItemID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "work_item_id is required"}}, nil
	}
	comments := resolveValue("comments", current, config)
	body := map[string]any{
		"requests": []any{
			map[string]any{
				"actionType": "Reject",
				"workitemId": workItemID,
				"comments":   comments,
			},
		},
	}
	result, err := client.post("/process/approvals", body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}
