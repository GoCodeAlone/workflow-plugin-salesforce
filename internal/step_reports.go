package internal

import (
	"context"
	"fmt"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// reportListStep implements step.salesforce_report_list
type reportListStep struct {
	name       string
	moduleName string
}

func newReportListStep(name string, config map[string]any) (*reportListStep, error) {
	return &reportListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *reportListStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	arr, obj, err := client.getArray("/analytics/reports")
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	if arr != nil {
		return &sdk.StepResult{Output: map[string]any{"reports": arr, "count": len(arr)}}, nil
	}
	return &sdk.StepResult{Output: obj}, nil
}

// reportDescribeStep implements step.salesforce_report_describe
type reportDescribeStep struct {
	name       string
	moduleName string
}

func newReportDescribeStep(name string, config map[string]any) (*reportDescribeStep, error) {
	return &reportDescribeStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *reportDescribeStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	reportID := resolveValue("report_id", current, config)
	if reportID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "report_id is required"}}, nil
	}
	path := fmt.Sprintf("/analytics/reports/%s/describe", reportID)
	result, err := client.get(path)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// reportRunStep implements step.salesforce_report_run
type reportRunStep struct {
	name       string
	moduleName string
}

func newReportRunStep(name string, config map[string]any) (*reportRunStep, error) {
	return &reportRunStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *reportRunStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	reportID := resolveValue("report_id", current, config)
	if reportID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "report_id is required"}}, nil
	}
	path := fmt.Sprintf("/analytics/reports/%s", reportID)
	includeDetails := resolveBool("include_details", current, config)
	if includeDetails {
		path += "?includeDetails=true"
	}
	result, err := client.post(path, nil)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// dashboardListStep implements step.salesforce_dashboard_list
type dashboardListStep struct {
	name       string
	moduleName string
}

func newDashboardListStep(name string, config map[string]any) (*dashboardListStep, error) {
	return &dashboardListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *dashboardListStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	arr, obj, err := client.getArray("/analytics/dashboards")
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	if arr != nil {
		return &sdk.StepResult{Output: map[string]any{"dashboards": arr, "count": len(arr)}}, nil
	}
	return &sdk.StepResult{Output: obj}, nil
}

// dashboardDescribeStep implements step.salesforce_dashboard_describe
type dashboardDescribeStep struct {
	name       string
	moduleName string
}

func newDashboardDescribeStep(name string, config map[string]any) (*dashboardDescribeStep, error) {
	return &dashboardDescribeStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *dashboardDescribeStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	dashboardID := resolveValue("dashboard_id", current, config)
	if dashboardID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "dashboard_id is required"}}, nil
	}
	path := fmt.Sprintf("/analytics/dashboards/%s/describe", dashboardID)
	result, err := client.get(path)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// dashboardRefreshStep implements step.salesforce_dashboard_refresh
type dashboardRefreshStep struct {
	name       string
	moduleName string
}

func newDashboardRefreshStep(name string, config map[string]any) (*dashboardRefreshStep, error) {
	return &dashboardRefreshStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *dashboardRefreshStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	dashboardID := resolveValue("dashboard_id", current, config)
	if dashboardID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "dashboard_id is required"}}, nil
	}
	path := fmt.Sprintf("/analytics/dashboards/%s", dashboardID)
	result, err := client.post(path, nil)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}
