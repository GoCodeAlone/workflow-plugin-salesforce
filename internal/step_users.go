package internal

import (
	"context"
	"fmt"
	"net/url"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// userGetStep implements step.salesforce_user_get
type userGetStep struct {
	name       string
	moduleName string
}

func newUserGetStep(name string, config map[string]any) (*userGetStep, error) {
	return &userGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *userGetStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	userID := resolveValue("user_id", current, config)
	if userID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "user_id is required"}}, nil
	}
	path := fmt.Sprintf("/sobjects/User/%s", userID)
	result, err := client.get(path)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// userListStep implements step.salesforce_user_list
type userListStep struct {
	name       string
	moduleName string
}

func newUserListStep(name string, config map[string]any) (*userListStep, error) {
	return &userListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *userListStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	soql := "SELECT Id, Name, Email, Username, IsActive FROM User"
	if filter := resolveValue("filter", current, config); filter != "" {
		soql += " WHERE " + filter
	}
	path := "/query?q=" + url.QueryEscape(soql)
	result, err := client.get(path)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	records, _ := result["records"].([]any)
	return &sdk.StepResult{Output: map[string]any{
		"users": records,
		"count": len(records),
	}}, nil
}

// userCreateStep implements step.salesforce_user_create
type userCreateStep struct {
	name       string
	moduleName string
}

func newUserCreateStep(name string, config map[string]any) (*userCreateStep, error) {
	return &userCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *userCreateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	fields := resolveMap("fields", current, config)
	if fields == nil {
		return &sdk.StepResult{Output: map[string]any{"error": "fields is required"}}, nil
	}
	result, err := client.post("/sobjects/User", fields)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// userUpdateStep implements step.salesforce_user_update
type userUpdateStep struct {
	name       string
	moduleName string
}

func newUserUpdateStep(name string, config map[string]any) (*userUpdateStep, error) {
	return &userUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *userUpdateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	userID := resolveValue("user_id", current, config)
	if userID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "user_id is required"}}, nil
	}
	fields := resolveMap("fields", current, config)
	if fields == nil {
		fields = map[string]any{}
	}
	path := fmt.Sprintf("/sobjects/User/%s", userID)
	result, err := client.patch(path, fields)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// identityGetStep implements step.salesforce_identity_get
type identityGetStep struct {
	name       string
	moduleName string
}

func newIdentityGetStep(name string, config map[string]any) (*identityGetStep, error) {
	return &identityGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *identityGetStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	// Identity endpoint is at the instance level, not versioned
	path := fmt.Sprintf("%s/services/oauth2/userinfo", client.instanceURL)
	result, err := client.get(path)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// orgLimitsStep implements step.salesforce_org_limits
type orgLimitsStep struct {
	name       string
	moduleName string
}

func newOrgLimitsStep(name string, config map[string]any) (*orgLimitsStep, error) {
	return &orgLimitsStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *orgLimitsStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	result, err := client.get("/limits")
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}
