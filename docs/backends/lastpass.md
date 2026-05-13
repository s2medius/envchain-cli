# LastPass Backend

The `lastpass` backend retrieves secrets from **LastPass Enterprise** using the LastPass API.

## Configuration

```yaml
version: 1
backends:
  - type: lastpass
    options:
      username: user@example.com
      api_key: your-enterprise-api-key
      api_url: https://lastpass.com/enterpriseapi.php  # optional
```

### Required Options

| Option     | Description                                      |
|------------|--------------------------------------------------|
| `username` | The LastPass Enterprise admin account username.  |
| `api_key`  | The LastPass Enterprise API key.                 |

### Optional Options

| Option    | Default                                     | Description                   |
|-----------|---------------------------------------------|-------------------------------|
| `api_url` | `https://lastpass.com/enterpriseapi.php`    | Override the LastPass API URL. |

## Key Format

Keys must be specified in the format `folder/name`, where:

- `folder` — the name of the LastPass shared folder.
- `name` — the name of the secret within that folder.

**Example:**

```yaml
vars:
  - key: DB_PASSWORD
    from: shared-infra/DB_PASSWORD
  - key: API_TOKEN
    from: shared-infra/API_TOKEN
```

## Notes

- This backend requires a **LastPass Enterprise** account with API access enabled.
- The `api_key` should be kept secret and ideally sourced from a secure environment variable.
- Ensure the admin user has access to the shared folders referenced in your configuration.
