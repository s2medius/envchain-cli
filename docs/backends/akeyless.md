# Akeyless Backend

The `akeyless` backend fetches secrets from [Akeyless](https://www.akeyless.io/), a unified secrets management platform.

## Configuration

```yaml
version: 1
backends:
  - type: akeyless
    options:
      token: "p-xxxxxxxxxxxxxxxx"
      gateway: "https://api.akeyless.io" # optional, defaults to public API
```

### Options

| Option    | Required | Default                      | Description                                      |
|-----------|----------|------------------------------|--------------------------------------------------|
| `token`   | yes      | —                            | Akeyless access token (API key or JWT)           |
| `gateway` | no       | `https://api.akeyless.io`    | URL of the Akeyless gateway (self-hosted or SaaS)|

## Key Format

Keys must be the **full secret path** as stored in Akeyless, e.g.:

```
/prod/myapp/db_password
```

## Example

```yaml
version: 1
backends:
  - type: akeyless
    options:
      token: "p-abc123"
vars:
  DB_PASSWORD: /prod/myapp/db_password
  API_KEY: /prod/myapp/api_key
```

Then run:

```bash
envchain run --config envchain.yaml -- ./server
```

## Notes

- The token must have **read** permission on all referenced secret paths.
- For self-hosted deployments, set `gateway` to your private gateway URL.
- Akeyless supports many authentication methods; use the API key token obtained from the Akeyless console or CLI (`akeyless create-api-key`).
