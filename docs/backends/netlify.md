# Netlify Backend

The `netlify` backend retrieves environment variables from [Netlify](https://www.netlify.com/) site environment configuration using the Netlify API.

## Configuration

```yaml
version: 1
backends:
  - name: netlify-prod
    type: netlify
    options:
      token: "${NETLIFY_TOKEN}"
      site_id: "your-site-id"
      api_url: "https://api.netlify.com/api/v1"  # optional
```

## Options

| Option    | Required | Default                                   | Description                          |
|-----------|----------|-------------------------------------------|--------------------------------------|
| `token`   | yes      | —                                         | Netlify personal access token        |
| `site_id` | yes      | —                                         | The Netlify site ID                  |
| `api_url` | no       | `https://api.netlify.com/api/v1`          | Base URL for the Netlify API         |

## Authentication

Generate a personal access token from the [Netlify user settings](https://app.netlify.com/user/applications#personal-access-tokens) page.

It is recommended to pass the token via an environment variable rather than hardcoding it in the config file:

```yaml
options:
  token: "${NETLIFY_AUTH_TOKEN}"
```

## Key Format

Keys correspond directly to environment variable names as defined in your Netlify site settings.

When multiple contexts are defined for a variable, `envchain-cli` will prefer the `all` or `production` context value.

## Example Usage

```yaml
version: 1
backends:
  - name: netlify
    type: netlify
    options:
      token: "${NETLIFY_AUTH_TOKEN}"
      site_id: "abc123def456"

variables:
  - key: DATABASE_URL
    backend: netlify
  - key: API_SECRET
    backend: netlify
```

```sh
envchain run --config envchain.yaml -- node server.js
```
