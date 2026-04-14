### Test Assertions with Testify

This project uses [testify](https://github.com/stretchr/testify) for all test assertions.
Testify provides two main assertion packages:

#### `assert` Package

Use `assert` when you want to continue executing the test even if an assertion fails.
This is useful when you want to check multiple conditions in a single test.

```go
import "github.com/stretchr/testify/assert"

func TestMultipleAssertions(t *testing.T) {
    result := SomeFunction()

    assert.NotNil(t, result)
    assert.Equal(t, expectedValue, result.Value)
    assert.Contains(t, result.Message, "expected text")
    // All assertions will be checked even if one fails
}
```

#### `require` Package

Use `require` when you want to stop the test immediately if an assertion fails.
This is useful when subsequent code depends on the assertion passing.

```go
import "github.com/stretchr/testify/require"

func TestWithDependency(t *testing.T) {
    result, err := SomeFunction()
    require.NoError(t, err)  // Test stops here if error occurs

    // This code only runs if the above assertion passes
    require.Equal(t, expectedValue, result.Value)
}
```

#### Common Assertions

**Equality:**

- `assert.Equal(t, expected, actual)` / `require.Equal(t, expected, actual)`
- `assert.NotEqual(t, expected, actual)` / `require.NotEqual(t, expected, actual)`

**Nil Checks:**

- `assert.Nil(t, value)` / `require.Nil(t, value)`
- `assert.NotNil(t, value)` / `require.NotNil(t, value)`

**Error Checks:**

- `assert.NoError(t, err)` / `require.NoError(t, err)`
- `assert.Error(t, err)` / `require.Error(t, err)`
- `assert.ErrorIs(t, err, target)` / `require.ErrorIs(t, err, target)`
- `assert.ErrorAs(t, err, target)` / `require.ErrorAs(t, err, target)`

**Boolean:**

- `assert.True(t, condition)` / `require.True(t, condition)`
- `assert.False(t, condition)` / `require.False(t, condition)`

**Contains/Subset:**

- `assert.Contains(t, container, item)` / `require.Contains(t, container, item)`
- `assert.Subset(t, subset, list)` / `require.Subset(t, subset, list)`

**Length/Count:**

- `assert.Len(t, object, length)` / `require.Len(t, object, length)`
- `assert.Empty(t, object)` / `require.Empty(t, object)`
- `assert.NotEmpty(t, object)` / `require.NotEmpty(t, object)`

For a complete list of assertions, see the [testify documentation](https://github.com/stretchr/testify).
