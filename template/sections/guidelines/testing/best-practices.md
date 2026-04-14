### Test Best Practices

**Do:**

- Write tests for all exported functions
- Use [testify](https://github.com/stretchr/testify) for all assertions (`assert` and `require`)
- Use `require` for critical assertions that must pass for the test to continue
- Use `assert` for non-critical assertions where you want to check multiple conditions
- Use table-driven tests for multiple cases
- Test both success and error paths
- Use descriptive test names
- Keep tests simple and focused
- Use mockery-generated mocks for infrastructure dependencies
- Use `EXPECT()` for type-safe mock expectations
- Pass `t *testing.T` to mock constructors
- Use integration tests for end-to-end verification

**Don't:**

- Don't use standard library `t.Errorf` or `t.Fatalf` for assertions (use testify instead)
- Don't use `reflect.DeepEqual` directly (use `assert.Equal` or `require.Equal` instead)
- Don't test implementation details
- Don't write flaky tests
- Don't skip error handling in tests
- Don't use sleeps for timing (use channels or mocks)
- Don't test private functions directly (test through public API)
- Don't write tests that depend on external state
- Don't manually verify mock expectations (automatic with testify)
- Don't over-mock (only mock direct dependencies of the code under test)
