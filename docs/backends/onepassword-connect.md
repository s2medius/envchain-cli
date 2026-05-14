# 1Password Connect Backend

The `onepassword_connect` backend retrieves secrets from a self-hosted
[1Password Connect Server](https://developer.1password.com/docs/connect).

## Configuration

| Field      | Required | Description                                      |
|------------|----------|--------------------------------------------------|
| `type`     | yes      | Must be `onepassword_connect`                    |
| `url`      | yes      | Base URL of your Connect Server (no trailing `/`)|
| `token`    | yes      | Connect Server access token                      |
| `vault_id` | yes      | UUID of the vault to query                       |

## Key Format

Keys must be in the format `<item-title>/<field-label>`, for example:

```
myapp/API_KEY
database/PASSWORD
```

Field label matching is case-insensitive.

## Example Config

```yaml
version: 1
backends:
  - name: op-connect
    type: onepassword_connect
    url: https://op-connect.internal
    token: ${OP_CONNECT_TOKEN}
    vault_id: aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee

variables:
  DATABASE_URL: op-connect:myapp/DATABASE_URL
  API_KEY: op-connect:myapp/API_KEY
```

## Running a Connect Server

See the [official documentation](https://developer.1password.com/docs/connect/get-started)
for instructions on deploying a Connect Server with Docker or Kubernetes.

## Security Notes

- Store your Connect token in an environment variable, not in the config file.
- Restrict the token's vault access to only the vaults your application needs.
- Use TLS for all Connect Server endpoints in production.
