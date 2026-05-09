# Infisical Backend

The `infisical` backend fetches secrets from [Infisical](https://infisical.com), an open-source secret management platform.

## Configuration

```yaml
version: 1
backends:
  - type: infisical
    options:
      token: "your-service-token"
      project_id: "your-project-id"
      environment: "prod"        # optional, defaults to "dev"
      base_url: "https://app.infisical.com"  # optional, for self-hosted instances

vars:
  DATABASE_URL: DATABASE_URL
  API_KEY: API_KEY
```

## Options

| Option        | Required | Default                        | Description                                      |
|---------------|----------|--------------------------------|--------------------------------------------------|
| `token`       | Yes      | —                              | Infisical service token for authentication       |
| `project_id`  | Yes      | —                              | Workspace / project ID in Infisical              |
| `environment` | No       | `dev`                          | Target environment (e.g. `dev`, `staging`, `prod`) |
| `base_url`    | No       | `https://app.infisical.com`    | Base URL for self-hosted Infisical instances     |

## Authentication

Generate a **Service Token** in the Infisical dashboard under **Project Settings → Service Tokens**.
Grant it read access to the secrets and environment you need.

## Example

```yaml
version: 1
backends:
  - type: infisical
    options:
      token: "st.abc123"
      project_id: "64f1a2b3c4d5e6f7a8b9c0d1"
      environment: "prod"

vars:
  DB_PASSWORD: DB_PASSWORD
  STRIPE_SECRET_KEY: STRIPE_SECRET_KEY
```

```sh
envchain-cli run --config envchain.yaml -- node server.js
```

## Notes

- Secret names in `vars` must match the secret name exactly as stored in Infisical.
- For self-hosted deployments, set `base_url` to your instance URL (e.g. `https://secrets.example.com`).
