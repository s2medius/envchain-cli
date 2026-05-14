# Keeper Secrets Manager Backend

The `keeper` backend retrieves secrets from [Keeper Secrets Manager](https://docs.keeper.io/secrets-manager/) via its REST gateway.

## Configuration

```yaml
version: 1
backends:
  - type: keeper
    options:
      token: "<your-access-token>"
      base_url: "https://keepersecurity.com/api/rest/sm/v1"  # optional
```

### Required Options

| Option  | Description                                  |
|---------|----------------------------------------------|
| `token` | Keeper Secrets Manager access token          |

### Optional Options

| Option     | Default                                              | Description                   |
|------------|------------------------------------------------------|-------------------------------|
| `base_url` | `https://keepersecurity.com/api/rest/sm/v1`          | Override the API base URL     |

## Key Format

Keys follow the format `<record_uid>/<field_type>`.

- `record_uid` — the UID of the Keeper record
- `field_type` — the field type to retrieve (e.g. `password`, `login`, `url`)

### Example

```yaml
variables:
  DB_PASSWORD: "AbCdEfGh12345678/password"
  DB_USER: "AbCdEfGh12345678/login"
```

## Authentication

The backend uses a Bearer token passed via the `Authorization` HTTP header. Generate an access token from the Keeper Secrets Manager console or via the KSM CLI.

## Notes

- Only the first value of a multi-value field is returned.
- The backend does **not** cache responses between calls.
