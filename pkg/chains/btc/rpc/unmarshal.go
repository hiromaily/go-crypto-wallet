package rpc

import (
	"encoding/json"
	"strings"
)

// FlexibleLabels handles both string arrays (BTC) and object arrays (BCH) for the labels field.
// BTC format: ["label1", "label2"]
// BCH format: [{"name": "label1", "purpose": "..."}, {"name": "label2", "purpose": "..."}]
type FlexibleLabels []string

// UnmarshalJSON implements custom unmarshaling for the labels field.
func (f *FlexibleLabels) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as string array (BTC format)
	var stringLabels []string
	if err := json.Unmarshal(data, &stringLabels); err == nil {
		*f = FlexibleLabels(stringLabels)
		return nil
	}

	// Try to unmarshal as object array (BCH format)
	var objectLabels []struct {
		Name    string `json:"name"`
		Purpose string `json:"purpose"`
	}
	if err := json.Unmarshal(data, &objectLabels); err == nil {
		labels := make([]string, len(objectLabels))
		for i, obj := range objectLabels {
			labels[i] = obj.Name
		}
		*f = FlexibleLabels(labels)
		return nil
	}

	// If both fail (e.g. null), return empty
	*f = FlexibleLabels([]string{})
	return nil
}

// WarningsField handles both string (Bitcoin Core < 29.2) and array (Bitcoin Core >= 29.2) formats
// for the warnings field in network info and blockchain info responses.
type WarningsField struct {
	Value []string
}

// UnmarshalJSON implements custom unmarshaling for the warnings field.
func (w *WarningsField) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as array first
	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		w.Value = arr
		return nil
	}

	// Fall back to string format
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	if str != "" {
		w.Value = []string{str}
	} else {
		w.Value = []string{}
	}
	return nil
}

// MarshalJSON marshals WarningsField as a JSON array of strings.
func (w WarningsField) MarshalJSON() ([]byte, error) {
	v := w.Value
	if v == nil {
		v = []string{}
	}
	return json.Marshal(v)
}

// String returns a concatenated string of all warnings.
func (w *WarningsField) String() string {
	if len(w.Value) == 0 {
		return ""
	}
	return strings.Join(w.Value, ", ")
}
