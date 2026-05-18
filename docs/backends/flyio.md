# Fly.io Backend

The `flyio` backend fetches secrets from the [Fly.io Secrets API](https://fly.io/docs/reference/secrets/).

## Configuration

| Option    | Required | Default                  | Description                        |
|-----------|----------|--------------------------|------------------------------------|
| `token`   | ✅        | —                        | Fly.io API token                   |
| `app_id`  | ✅        | —                        | The Fly.io application name/ID     |
| `api_url` | ❌        | `https://api.fly.io`     | Override the API base URL          |

## Example

```yaml
version: 1
backends:
  - name: fly
    type: flyio
    options:
      token: "${FLY_API_TOKEN}"
      app_id: "my-production-app"

vars:
  DATABASE_URL:
    backend: fly
    key: DATABASE_URL
  API_KEY:
    backend: fly
    key: API_KEY
```

## Authentication

Generate a token via the Fly.io CLI:

```bash
fly auth token
```

Or create a deploy token scoped to a specific app:

```bash
fly tokens create deploy -a my-production-app
```

## Notes

- Secrets are fetched in a single API call and cached for the duration of resolution.
- The key name must match the exact secret name set in Fly.io.
- Only secrets visible to the provided token will be accessible.
