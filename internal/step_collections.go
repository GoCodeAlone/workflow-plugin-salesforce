package internal

import (
	"context"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// collectionInsertStep implements step.salesforce_collection_insert
type collectionInsertStep struct {
	name       string
	moduleName string
}

func newCollectionInsertStep(name string, config map[string]any) (*collectionInsertStep, error) {
	return &collectionInsertStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *collectionInsertStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	records := resolveAnySlice("records", current, config)
	if len(records) == 0 {
		return &sdk.StepResult{Output: map[string]any{"error": "records is required"}}, nil
	}
	allOrNone := resolveBool("all_or_none", current, config)
	body := map[string]any{
		"records":   records,
		"allOrNone": allOrNone,
	}
	arr, obj, err := client.postArray("/composite/sobjects", body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	if arr != nil {
		return &sdk.StepResult{Output: map[string]any{"results": arr, "count": len(arr)}}, nil
	}
	return &sdk.StepResult{Output: obj}, nil
}

// collectionUpdateStep implements step.salesforce_collection_update
type collectionUpdateStep struct {
	name       string
	moduleName string
}

func newCollectionUpdateStep(name string, config map[string]any) (*collectionUpdateStep, error) {
	return &collectionUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *collectionUpdateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	records := resolveAnySlice("records", current, config)
	if len(records) == 0 {
		return &sdk.StepResult{Output: map[string]any{"error": "records is required"}}, nil
	}
	allOrNone := resolveBool("all_or_none", current, config)
	body := map[string]any{
		"records":   records,
		"allOrNone": allOrNone,
	}
	result, err := client.patch("/composite/sobjects", body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// collectionUpsertStep implements step.salesforce_collection_upsert
type collectionUpsertStep struct {
	name       string
	moduleName string
}

func newCollectionUpsertStep(name string, config map[string]any) (*collectionUpsertStep, error) {
	return &collectionUpsertStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *collectionUpsertStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	sObjectType := resolveValue("sobject_type", current, config)
	externalIDField := resolveValue("external_id_field", current, config)
	records := resolveAnySlice("records", current, config)
	if sObjectType == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "sobject_type is required"}}, nil
	}
	if externalIDField == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "external_id_field is required"}}, nil
	}
	if len(records) == 0 {
		return &sdk.StepResult{Output: map[string]any{"error": "records is required"}}, nil
	}
	allOrNone := resolveBool("all_or_none", current, config)
	body := map[string]any{
		"records":   records,
		"allOrNone": allOrNone,
	}
	path := "/composite/sobjects/" + sObjectType + "/" + externalIDField
	result, err := client.patch(path, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// collectionDeleteStep implements step.salesforce_collection_delete
type collectionDeleteStep struct {
	name       string
	moduleName string
}

func newCollectionDeleteStep(name string, config map[string]any) (*collectionDeleteStep, error) {
	return &collectionDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *collectionDeleteStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	ids := resolveStringSlice("ids", current, config)
	if len(ids) == 0 {
		return &sdk.StepResult{Output: map[string]any{"error": "ids is required"}}, nil
	}
	allOrNone := resolveBool("all_or_none", current, config)
	idParam := ""
	for i, id := range ids {
		if i > 0 {
			idParam += ","
		}
		idParam += id
	}
	allOrNoneParam := "false"
	if allOrNone {
		allOrNoneParam = "true"
	}
	path := "/composite/sobjects?ids=" + idParam + "&allOrNone=" + allOrNoneParam
	result, err := client.delete(path)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// resolveAnySlice resolves a []any from current or config.
func resolveAnySlice(key string, current, config map[string]any) []any {
	for _, m := range []map[string]any{current, config} {
		switch v := m[key].(type) {
		case []any:
			return v
		case []map[string]any:
			result := make([]any, len(v))
			for i, item := range v {
				result[i] = item
			}
			return result
		}
	}
	return nil
}
