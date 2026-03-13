package internal

import "fmt"

func getModuleName(config map[string]any) string {
	if v, ok := config["module"].(string); ok && v != "" {
		return v
	}
	return "salesforce"
}

func resolveValue(key string, current, config map[string]any) string {
	if v, ok := current[key].(string); ok && v != "" {
		return v
	}
	if v, ok := config[key].(string); ok && v != "" {
		return v
	}
	return ""
}

func resolveInt64(key string, current, config map[string]any) int64 {
	if v := toInt64(current[key]); v != 0 {
		return v
	}
	return toInt64(config[key])
}

func resolveFloat64(key string, current, config map[string]any) float64 {
	if v := toFloat64(current[key]); v != 0 {
		return v
	}
	return toFloat64(config[key])
}

func resolveStringSlice(key string, current, config map[string]any) []string {
	for _, m := range []map[string]any{current, config} {
		switch v := m[key].(type) {
		case []string:
			return v
		case []any:
			result := make([]string, 0, len(v))
			for _, item := range v {
				if s, ok := item.(string); ok {
					result = append(result, s)
				}
			}
			return result
		}
	}
	return nil
}

func resolveMap(key string, current, config map[string]any) map[string]any {
	if v, ok := current[key].(map[string]any); ok {
		return v
	}
	if v, ok := config[key].(map[string]any); ok {
		return v
	}
	return nil
}

func resolveBool(key string, current, config map[string]any) bool {
	for _, m := range []map[string]any{current, config} {
		switch v := m[key].(type) {
		case bool:
			return v
		case string:
			return v == "true" || v == "1" || v == "yes"
		}
	}
	return false
}

func resolveInt(key string, current, config map[string]any) int {
	return int(resolveInt64(key, current, config))
}

func toInt64(v any) int64 {
	switch t := v.(type) {
	case int64:
		return t
	case int:
		return int64(t)
	case int32:
		return int64(t)
	case float64:
		return int64(t)
	case float32:
		return int64(t)
	case string:
		var n int64
		fmt.Sscanf(t, "%d", &n)
		return n
	}
	return 0
}

func toFloat64(v any) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case float32:
		return float64(t)
	case int64:
		return float64(t)
	case int:
		return float64(t)
	case string:
		var f float64
		fmt.Sscanf(t, "%f", &f)
		return f
	}
	return 0
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func derefBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// suppress unused warnings for helpers used only in some files
var (
	_ = resolveFloat64
	_ = resolveInt
	_ = resolveStringSlice
	_ = derefInt
	_ = derefBool
	_ = resolveMap
	_ = resolveBool
)
