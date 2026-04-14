### Dependency Management

**Go Modules:**

- Use `go mod tidy` to organize dependencies
- Run security scans (`govulncheck`)
- Keep dependencies up-to-date while maintaining stability

**Commands:**

- `make tidy`: Organize dependencies and clean up `go.mod`
- `make go-check-vuln`: Run security vulnerability scan (govulncheck)

**Best Practices:**

- Review dependency changes carefully
- Test thoroughly after dependency updates
- Document breaking changes in dependencies
- Consider security implications of new dependencies
