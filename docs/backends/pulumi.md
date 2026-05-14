# Pulumi ESC Backend

The `pulumi` backend retrieves secrets from [Pulumi ESC](https://www.pulumi.com/product/esc/) (Environments, Secrets, and Configuration).

## Configuration

| Key           | Required | Default                      | Description                              |
|---------------|----------|------------------------------|------------------------------------------|
| `token`       | Yes      | —                            | Pulumi access token                      |
| `org`         | Yes      | —                            | Pulumi organization name                 |
| `environment` | Yes      | —                            | ESC environment name                     |
| `api_url`     | No       | `https://api.pulumi.com`     | Override the Pulumi API base URL         |

## Example Config

```yaml
version: 1
backends:
  - name: pulumi-prod
    type: pulumi
    config:
      token: ${PULUMI_ACCESS_TOKEN}
      org: my-org
      environment: prod

vars:
  - name: DATABASE_URL
    backend: pulumi-prod
    key: DATABASE_URL
  - name: API_KEY
    backend: pulumi-prod
    key: API_KEY
```

## Authentication

Generate a Pulumi access token at https://app.pulumi.com/account/tokens.

It is recommended to store the token in an environment variable and reference it via `${PULUMI_ACCESS_TOKEN}` in your config file rather than hardcoding it.

## Key Format

Keys correspond to the top-level value names defined in your Pulumi ESC environment. For example, if your environment defines:

```yaml
values:
  DATABASE_URL: postgres://...
  API_KEY: abc123
```

You would reference them as `DATABASE_URL` and `API_KEY` respectively.

## Notes

- The backend opens the environment using the Pulumi ESC REST API and reads values from the `values` map.
- Nested keys are not currently supported; only top-level values are accessible.
