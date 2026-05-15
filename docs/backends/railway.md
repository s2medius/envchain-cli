# Railway Backend

The `railway` backend retrieves environment variables from [Railway](https://railway.app) projects using the Railway GraphQL API.

## Configuration

```yaml
version: 1
backends:
  - name: railway
    type: railway
    options:
      token: "${RAILWAY_API_TOKEN}"
      project_id: "your-project-id"
      environment_id: "your-environment-id"
```

## Options

| Option           | Required | Default                                        | Description                              |
|------------------|----------|------------------------------------------------|------------------------------------------|
| `token`          | Yes      | —                                              | Railway API token                        |
| `project_id`     | Yes      | —                                              | Railway project ID                       |
| `environment_id` | Yes      | —                                              | Railway environment ID                   |
| `api_url`        | No       | `https://backboard.railway.app/graphql/v2`     | Override the Railway GraphQL API URL     |

## Finding Your IDs

You can find your **project ID** and **environment ID** in the Railway dashboard URL:

```
https://railway.app/project/<project_id>/environment/<environment_id>
```

Alternatively, use the Railway CLI:

```bash
railway status
```

## Obtaining an API Token

1. Go to [railway.app/account/tokens](https://railway.app/account/tokens)
2. Click **New Token**
3. Copy the generated token and store it securely

## Usage Example

```yaml
version: 1
backends:
  - name: prod-railway
    type: railway
    options:
      token: "${RAILWAY_TOKEN}"
      project_id: "abc-123"
      environment_id: "production"

vars:
  - key: DATABASE_URL
    backend: prod-railway
  - key: REDIS_URL
    backend: prod-railway
```

## Notes

- The backend fetches **all** variables for the given project/environment in a single API call and filters by key name locally.
- Railway API tokens are scoped to a team or personal account; ensure the token has read access to the target project.
