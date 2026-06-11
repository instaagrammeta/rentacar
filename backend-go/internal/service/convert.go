package service

import (
	"time"

	"github.com/instaagrammeta/rentacar/backend-go/internal/models"
)

// Helpers to read values from a JSON-decoded map (numbers arrive as float64)
// when importing a .rentacar project file.

func asStr(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func asStrPtr(m map[string]interface{}, key string) *string {
	if v, ok := m[key].(string); ok && v != "" {
		return &v
	}
	return nil
}

func asFloat(m map[string]interface{}, key string) float64 {
	switch v := m[key].(type) {
	case float64:
		return v
	case int:
		return float64(v)
	}
	return 0
}

func asInt(m map[string]interface{}, key string) int {
	switch v := m[key].(type) {
	case float64:
		return int(v)
	case int:
		return v
	}
	return 0
}

func asIntPtr(m map[string]interface{}, key string) *int {
	switch v := m[key].(type) {
	case float64:
		i := int(v)
		return &i
	case int:
		return &v
	}
	return nil
}

func asUint(m map[string]interface{}, key string) uint {
	switch v := m[key].(type) {
	case float64:
		return uint(v)
	case int:
		return uint(v)
	}
	return 0
}

func asUintPtr(m map[string]interface{}, key string) *uint {
	switch v := m[key].(type) {
	case float64:
		u := uint(v)
		return &u
	case int:
		u := uint(v)
		return &u
	}
	return nil
}

func asBool(m map[string]interface{}, key string) bool {
	if v, ok := m[key].(bool); ok {
		return v
	}
	return false
}

// asDatePtr parses a "YYYY-MM-DD" string into a *time.Time.
func asDatePtr(m map[string]interface{}, key string) *time.Time {
	if v, ok := m[key].(string); ok && v != "" {
		if t, err := time.Parse("2006-01-02", v[:min(10, len(v))]); err == nil {
			return &t
		}
	}
	return nil
}

// asDateVal parses a required date column.
func asDateVal(m map[string]interface{}, key string) time.Time {
	if p := asDatePtr(m, key); p != nil {
		return *p
	}
	return time.Now()
}

// asDateTime parses an ISO datetime column.
func asDateTime(m map[string]interface{}, key string) time.Time {
	if p := asDateTimePtr(m, key); p != nil {
		return *p
	}
	return time.Now()
}

// asDateTimePtr parses an optional ISO datetime column (nil when absent).
func asDateTimePtr(m map[string]interface{}, key string) *time.Time {
	if v, ok := m[key].(string); ok && v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			return &t
		}
		if t, err := time.Parse("2006-01-02T15:04:05", v); err == nil {
			return &t
		}
	}
	return nil
}

// base reconstructs a models.Base (id + timestamps) from an exported row.
func base(m map[string]interface{}) models.Base {
	b := models.Base{ID: asUint(m, "id")}
	if v, ok := m["created_at"].(string); ok && v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			b.CreatedAt = t
		}
	}
	if v, ok := m["updated_at"].(string); ok && v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			b.UpdatedAt = t
		}
	}
	return b
}
