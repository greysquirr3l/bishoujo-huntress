# API Reference

## Client Configuration

### huntress.New()

Creates a new Huntress API client with configuration options.

**Input:**
- `options ...ClientOption`: Configuration options

**Output:**
- `*Client`: Configured client instance

**Example:**

```go
client := huntress.New(
    huntress.WithCredentials("API_KEY", "API_SECRET"),
    huntress.WithTimeout(30 * time.Second),
)
```

## Organization Service

### List Organizations

**Method:** `client.Organization.List(ctx, params)`

**Input:**
- `ctx context.Context`: Request context
- `params *organization.ListParams`: Query parameters
  - `Page int`: Page number (default: 1)
  - `Limit int`: Items per page (max: 100)

**Output:**
- `[]*organization.Organization`: List of organizations
- `*Pagination`: Pagination metadata
- `error`: Error if any

**Example:**

```go
orgs, pagination, err := client.Organization.List(ctx, &organization.ListParams{
    Page:  1,
    Limit: 50,
})
```

// Continue for all public methods...
