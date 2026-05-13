# HCP Vault Secrets Backend

The `hashicorp` backend fetches secrets from [HashiCorp Cloud Platform (HCP) Vault Secrets](https://developer.hashicorp.com/hcp/docs/vault-secrets).

## Configuration

| Field        | Required | Description                                      |
|--------------|----------|--------------------------------------------------|
| `token`      | ✅        | HCP service principal token                      |
| `org_id`     | ✅        | HCP organization ID                              |
| `project_id` | ✅        | HCP project ID                                   |
| `app_name`   | ✅        | Name of the HCP Vault Secrets application        |
| `base_url`   | ❌        | Override API base URL (default: HCP production)  |

## Example

```yaml
version: 1
backends:
  - name: hcp
    type: hashicorp
    options:
      token: "${HCP_TOKEN}"
      org_id: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
      project_id: "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
      app_name: "my-app"

variables:
  DATABASE_URL:
    backend: hcp
    key: DATABASE_URL
  API_KEY:
    backend: hcp
    key: API_KEY
```

## Obtaining Credentials

1. Log in to [HCP Portal](https://portal.cloud.hashicorp.com).
2. Navigate to **Access Control (IAM) → Service Principals**.
3. Create a service principal and generate a client secret.
4. Exchange the client credentials for a bearer token via the HCP auth API.

## Notes

- Secrets must be stored under the specified **app** in HCP Vault Secrets.
- The `key` field in the variable mapping corresponds to the **secret name** in HCP.
- Only the latest static version of the secret is retrieved.
