# Bitwarden Backend

The Bitwarden backend retrieves secrets from a [Bitwarden](https://bitwarden.com) vault item using the [Bitwarden CLI REST server](https://bitwarden.com/help/cli/#serve) (`bw serve`).

Secrets are read from the **custom fields** of a single vault item identified by its item ID.

## Prerequisites

1. Install the [Bitwarden CLI](https://bitwarden.com/help/cli/).
2. Log in and unlock your vault:
   ```bash
   bw login
   bw unlock
   ```
3. Start the local REST server:
   ```bash
   bw serve --port 8087
   ```

## Configuration

```yaml
version: 1
backends:
  - name: bw
    type: bitwarden
    options:
      access_token: "<BW_SESSION_TOKEN>"   # required
      item_id: "<vault-item-uuid>"          # required
      base_url: "http://localhost:8087"     # optional, default shown

vars:
  - key: DB_PASSWORD
    backend: bw
    field: DB_PASSWORD
```

### Options

| Option         | Required | Default                   | Description                                      |
|----------------|----------|---------------------------|--------------------------------------------------|
| `access_token` | Yes      | —                         | Bitwarden session token (`BW_SESSION` env value) |
| `item_id`      | Yes      | —                         | UUID of the vault item containing the fields     |
| `base_url`     | No       | `http://localhost:8087`   | Base URL of the `bw serve` REST server           |

## How It Works

The backend calls `GET /object/item/<item_id>` on the local Bitwarden CLI server and maps each custom field `name` → `value`. When `envchain-cli` resolves a variable, it looks up the field by the `field` key specified in the config.

## Security Notes

- Keep the `bw serve` process local and never expose it on a public interface.
- Prefer passing `access_token` via an environment variable rather than hard-coding it in your config file.
- Add your config file to `.gitignore` to avoid committing secrets.
