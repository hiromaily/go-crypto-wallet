package postgres

import "fmt"

// interfaceToString safely converts interface{} to string.
// PostgreSQL SQLC generates enum columns as interface{} type.
func interfaceToString(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	if b, ok := v.([]byte); ok {
		return string(b)
	}
	return fmt.Sprintf("%v", v)
}
