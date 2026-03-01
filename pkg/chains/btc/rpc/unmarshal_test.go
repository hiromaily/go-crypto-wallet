package rpc

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── FlexibleLabels ───────────────────────────────────────────────────────────

func TestFlexibleLabels_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input string
		want  FlexibleLabels
	}{
		{
			name:  "BTC string array",
			input: `["label1", "label2"]`,
			want:  FlexibleLabels{"label1", "label2"},
		},
		{
			name:  "BCH object array",
			input: `[{"name": "label1", "purpose": "receive"}, {"name": "label2", "purpose": "change"}]`,
			want:  FlexibleLabels{"label1", "label2"},
		},
		{
			name:  "empty array",
			input: `[]`,
			want:  FlexibleLabels{},
		},
		{
			name:  "null value",
			input: `null`,
			want:  nil, // json.Unmarshal of null into []string yields nil
		},
		{
			name:  "single string element",
			input: `["only-label"]`,
			want:  FlexibleLabels{"only-label"},
		},
		{
			name:  "BCH single object",
			input: `[{"name": "my-account", "purpose": "receive"}]`,
			want:  FlexibleLabels{"my-account"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var fl FlexibleLabels
			err := json.Unmarshal([]byte(tt.input), &fl)
			require.NoError(t, err)
			assert.Equal(t, tt.want, fl)
		})
	}
}

// ─── WarningsField ────────────────────────────────────────────────────────────

func TestWarningsField_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		{
			name:  "array format (Bitcoin Core >= 29.2)",
			input: `["warning one", "warning two"]`,
			want:  []string{"warning one", "warning two"},
		},
		{
			name:  "string format (Bitcoin Core < 29.2) with content",
			input: `"this is a warning"`,
			want:  []string{"this is a warning"},
		},
		{
			name:  "empty string format",
			input: `""`,
			want:  []string{},
		},
		{
			name:  "empty array format",
			input: `[]`,
			want:  []string{},
		},
		{
			name:  "single warning in array",
			input: `["only warning"]`,
			want:  []string{"only warning"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var wf WarningsField
			err := json.Unmarshal([]byte(tt.input), &wf)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, wf.Value)
		})
	}
}

func TestWarningsField_String(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		value []string
		want  string
	}{
		{
			name:  "empty",
			value: []string{},
			want:  "",
		},
		{
			name:  "single warning",
			value: []string{"warning one"},
			want:  "warning one",
		},
		{
			name:  "multiple warnings",
			value: []string{"warning one", "warning two"},
			want:  "warning one, warning two",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			wf := &WarningsField{Value: tt.value}
			assert.Equal(t, tt.want, wf.String())
		})
	}
}
