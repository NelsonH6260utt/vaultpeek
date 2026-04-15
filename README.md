# vaultpeek

A CLI tool for browsing and diffing HashiCorp Vault secret paths across environments.

---

## Installation

```bash
go install github.com/yourname/vaultpeek@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/vaultpeek.git
cd vaultpeek
go build -o vaultpeek .
```

---

## Usage

Set your Vault address and token, then start browsing:

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.yourtoken"

# Browse secrets at a given path
vaultpeek browse secret/myapp

# Diff secrets between two environments
vaultpeek diff secret/myapp/staging secret/myapp/production
```

### Example Output

```
[~] DB_HOST       staging.db.internal  →  prod.db.internal
[+] FEATURE_FLAG  (missing)            →  "true"
[=] APP_PORT      8080                 ==  8080
```

---

## Configuration

| Flag            | Env Var       | Description                  |
|-----------------|---------------|------------------------------|
| `--addr`        | `VAULT_ADDR`  | Vault server address         |
| `--token`       | `VAULT_TOKEN` | Vault authentication token   |
| `--output`      | —             | Output format: `text`, `json`|

---

## Requirements

- Go 1.21+
- HashiCorp Vault 1.x

---

## License

MIT © 2024 yourname