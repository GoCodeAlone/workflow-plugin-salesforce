package internal

import (
	"context"
	"fmt"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// chatterPostStep implements step.salesforce_chatter_post
type chatterPostStep struct {
	name       string
	moduleName string
}

func newChatterPostStep(name string, config map[string]any) (*chatterPostStep, error) {
	return &chatterPostStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *chatterPostStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	subjectID := resolveValue("subject_id", current, config)
	text := resolveValue("text", current, config)
	if subjectID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "subject_id is required"}}, nil
	}
	if text == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "text is required"}}, nil
	}
	body := map[string]any{
		"body": map[string]any{
			"messageSegments": []any{
				map[string]any{"type": "Text", "text": text},
			},
		},
		"feedElementType": "FeedItem",
		"subjectId":       subjectID,
	}
	result, err := client.post("/chatter/feed-elements", body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// chatterCommentStep implements step.salesforce_chatter_comment
type chatterCommentStep struct {
	name       string
	moduleName string
}

func newChatterCommentStep(name string, config map[string]any) (*chatterCommentStep, error) {
	return &chatterCommentStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *chatterCommentStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	feedElementID := resolveValue("feed_element_id", current, config)
	text := resolveValue("text", current, config)
	if feedElementID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "feed_element_id is required"}}, nil
	}
	if text == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "text is required"}}, nil
	}
	body := map[string]any{
		"body": map[string]any{
			"messageSegments": []any{
				map[string]any{"type": "Text", "text": text},
			},
		},
	}
	path := fmt.Sprintf("/chatter/feed-elements/%s/capabilities/comments/items", feedElementID)
	result, err := client.post(path, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// chatterLikeStep implements step.salesforce_chatter_like
type chatterLikeStep struct {
	name       string
	moduleName string
}

func newChatterLikeStep(name string, config map[string]any) (*chatterLikeStep, error) {
	return &chatterLikeStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *chatterLikeStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	feedElementID := resolveValue("feed_element_id", current, config)
	if feedElementID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "feed_element_id is required"}}, nil
	}
	path := fmt.Sprintf("/chatter/feed-elements/%s/capabilities/chatter-likes/items", feedElementID)
	result, err := client.post(path, nil)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// chatterFeedListStep implements step.salesforce_chatter_feed_list
type chatterFeedListStep struct {
	name       string
	moduleName string
}

func newChatterFeedListStep(name string, config map[string]any) (*chatterFeedListStep, error) {
	return &chatterFeedListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *chatterFeedListStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "salesforce client not found: " + s.moduleName}}, nil
	}
	feedType := resolveValue("feed_type", current, config)
	if feedType == "" {
		feedType = "news"
	}
	path := fmt.Sprintf("/chatter/feeds/%s/me/feed-elements", feedType)
	result, err := client.get(path)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}
