### Naming Conventions

**Packages:**

- Use lowercase package names
- Avoid underscores in package names
- Keep package names short and descriptive

**Exported Symbols:**

- Start with uppercase letter
- Use camelCase (e.g., `GetAccountKey`, `CreateTransaction`)

**Unexported Symbols:**

- Start with lowercase letter
- Use camelCase (e.g., `calculateFee`, `validateAddress`)

**Constants:**

- Use camelCase for exported constants (e.g., `MaxRetries`)
- Use camelCase for unexported constants (e.g., `defaultTimeout`)

**Interfaces:**

- Name interfaces after the behavior they represent (e.g., `Validator`, `Repository`)
- For single-method interfaces, use the method name + "er" suffix (e.g., `Reader`, `Writer`)
