# 1Password Backend

The `1password` backend retrieves secrets from [1Password](https://1password.com) using the official [`op` CLI](https://developer.1password.com/docs/cli).

## Prerequisites

- `op` CLI installed and available on `$PATH`
- Signed in to your 1Password account (`op signin`)

## Configuration

```yaml
version: 1
backends:
  - name: op
    type: 1password
    options:
      vault: "dev"          # optional: restrict to a specific vault
      account: "myaccount"  # optional: specify account shorthand

vars:
  DB_PASSWORD:
    backend: op
    key: "db-credentials/password"   # format: <item>/<field>
  API_KEY:
    backend: op
    key: "api-service/credential"
```

## Key Format

Keys must follow the `item/field` format:

| Part   | Description                          |
|--------|--------------------------------------|
| `item` | The 1Password item name or UUID      |
| `field`| The field label within that item     |

## Options

| Option    | Required | Description                                  |
|-----------|----------|----------------------------------------------|
| `vault`   | No       | Vault name or UUID to scope lookups          |
| `account` | No       | Account shorthand (useful for multiple accounts) |

## Notes

- The backend shells out to `op` for each secret retrieval.
- Ensure the `op` session is active before running `envchain-cli`.
- Use `op run` as an alternative if you prefer native 1Password process injection.
