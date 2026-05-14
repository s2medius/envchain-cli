# Conjur Backend

The `conjur` backend retrieves secrets from [CyberArk Conjur](https://www.conjur.org/), an open-source secrets management solution.

## Configuration

| Option    | Required | Description                                      |
|-----------|----------|--------------------------------------------------|
| `address` | Yes      | Base URL of the Conjur appliance (e.g. `https://conjur.example.com`) |
| `account` | Yes      | Conjur organization account name                 |
| `token`   | Yes      | Short-lived Conjur API access token (base64)     |

## Example

```yaml
version: 1
backends:
  - name: conjur-prod
    type: conjur
    options:
      address: https://conjur.example.com
      account: myorg
      token: "${CONJUR_TOKEN}"

variables:
  DB_PASSWORD:
    backend: conjur-prod
    key: myapp/db/password
  API_KEY:
    backend: conjur-prod
    key: myapp/api/key
```

## Key Format

Keys follow the Conjur variable ID format: `<policy-branch>/<variable-name>`.

For example, if your variable is defined under the `myapp` policy as `db/password`, the key would be:

```
myapp/db/password
```

## Token Acquisition

Conjur access tokens are short-lived (default 8 minutes). You can obtain one using the Conjur CLI:

```bash
export CONJUR_TOKEN=$(conjur authenticate -H | base64)
```

Or via the REST API:

```bash
curl -s -X POST https://conjur.example.com/authn/myorg/<login>/authenticate \
  --data-binary @api_key.txt | base64
```

## Notes

- Ensure the identity associated with the token has `execute` privilege on the target variables.
- The `address` value must not include a trailing slash.
