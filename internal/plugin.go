// Package internal implements the workflow-plugin-salesforce plugin.
package internal

import (
	"fmt"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// Version is set at build time via -ldflags
// "-X github.com/GoCodeAlone/workflow-plugin-salesforce/internal.Version=X.Y.Z"
var Version = "dev"

// salesforcePlugin implements sdk.PluginProvider, sdk.ModuleProvider, and sdk.StepProvider.
type salesforcePlugin struct{}

// NewSalesforcePlugin returns a new salesforcePlugin instance.
func NewSalesforcePlugin() sdk.PluginProvider {
	return &salesforcePlugin{}
}

// Manifest returns plugin metadata.
func (p *salesforcePlugin) Manifest() sdk.PluginManifest {
	return sdk.PluginManifest{
		Name:        "workflow-plugin-salesforce",
		Version:     Version,
		Author:      "GoCodeAlone",
		Description: "Salesforce CRM platform plugin (~75 step types across all Salesforce REST APIs)",
	}
}

// ModuleTypes returns the module type names this plugin provides.
func (p *salesforcePlugin) ModuleTypes() []string {
	return []string{"salesforce.provider"}
}

// CreateModule creates a module instance of the given type.
func (p *salesforcePlugin) CreateModule(typeName, name string, config map[string]any) (sdk.ModuleInstance, error) {
	switch typeName {
	case "salesforce.provider":
		m, err := newSalesforceModule(name, config)
		if err != nil {
			return nil, err
		}
		return m, nil
	default:
		return nil, fmt.Errorf("salesforce plugin: unknown module type %q", typeName)
	}
}

// StepTypes returns the step type names this plugin provides.
func (p *salesforcePlugin) StepTypes() []string {
	return allStepTypes()
}

// CreateStep creates a step instance of the given type.
func (p *salesforcePlugin) CreateStep(typeName, name string, config map[string]any) (sdk.StepInstance, error) {
	return createStep(typeName, name, config)
}
