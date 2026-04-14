### Table-Driven Tests

Use table-driven tests for multiple test cases with testify assertions:

```go
import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   InputType
        want    OutputType
        wantErr bool
    }{
        {
            name:    "valid case",
            input:   validInput,
            want:    expectedOutput,
            wantErr: false,
        },
        {
            name:    "error case",
            input:   invalidInput,
            want:    nil,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Function(tt.input)
            if tt.wantErr {
                require.Error(t, err, "Function() should return error for input %v", tt.input)
                assert.Nil(t, got, "Function() should return nil result on error")
            } else {
                require.NoError(t, err, "Function() should not return error for input %v", tt.input)
                assert.Equal(t, tt.want, got, "Function() should return expected result")
            }
        })
    }
}
```
