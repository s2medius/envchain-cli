# envchain-cli

> CLI tool to securely inject environment variables from multiple secret backends into processes

---

## Installation

```bash
go install github.com/yourusername/envchain-cli@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/envchain-cli/releases).

---

## Usage

Inject secrets from a backend into any process by prefixing your command with `envchain`:

```bash
envchain --backend vault --path secret/myapp -- ./myapp serve
```

### Supported Backends

- **HashiCorp Vault** (`vault`)
- **AWS Secrets Manager** (`aws-ssm`)
- **1Password** (`onepassword`)
- **dotenv files** (`dotenv`)

### Example with multiple backends

```bash
envchain \
  --backend vault --path secret/db \
  --backend aws-ssm --path /prod/api \
  -- python manage.py runserver
```

Secrets are resolved at runtime, merged in order, and injected as environment variables into the child process. They are never written to disk or shell history.

### Configuration file

You can define backends in a config file (`envchain.yaml`):

```yaml
backends:
  - type: vault
    path: secret/myapp
  - type: dotenv
    path: .env.local
```

Then run:

```bash
envchain --config envchain.yaml -- ./myapp
```

---

## Contributing

Pull requests and issues are welcome. Please open an issue before submitting large changes.

---

## License

[MIT](LICENSE) © 2024 Your Name