# GitHub Actions Secrets Backend

The `github` backend fetches secrets from **GitHub Actions repository secrets** using the GitHub REST API.

> **Note:** GitHub's API only exposes secret *metadata* for most endpoints. The `value` field is available only in specific enterprise or self-hosted configurations. This backend is best suited for custom GitHub Enterprise setups or mock/test environments that return secret values directly.

## Configuration

```yaml
version: 1
backends:
  - name: gh
    type: github
    options:
      token: "ghp_yourPersonalAccessToken"
      owner: "my-org"
      repo: "my-repo"
      base_url: "https://api.github.com"  # optional, defaults to https://api.github.com
```

## Options

| Option     | Required | Default                      | Description                          |
|------------|----------|------------------------------|--------------------------------------|
| `token`    | Yes      | —                            | GitHub personal access token (PAT)   |
| `owner`    | Yes      | —                            | GitHub organization or user name     |
| `repo`     | Yes      | —                            | Repository name                      |
| `base_url` | No       | `https://api.github.com`     | Override for GitHub Enterprise Server |

## Required Token Scopes

Your PAT must have the following scopes:

- `repo` — for private repositories
- `public_repo` — for public repositories

## Example Usage

```yaml
version: 1
backends:
  - name: gh
    type: github
    options:
      token: "${GITHUB_TOKEN}"
      owner: "acme-corp"
      repo: "api-service"

variables:
  DATABASE_URL:
    backend: gh
    key: DATABASE_URL
  API_KEY:
    backend: gh
    key: THIRD_PARTY_API_KEY
```

Then run:

```bash
envchain run --config envchain.yaml -- ./start-server
```

## GitHub Enterprise Server

For GitHub Enterprise, set `base_url` to your instance URL:

```yaml
options:
  base_url: "https://github.example.com/api/v3"
```
