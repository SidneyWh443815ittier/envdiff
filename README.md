# envdiff

> Compare `.env` files across environments and flag missing, extra, or mismatched keys — with optional CI integration.

---

## Installation

```bash
go install github.com/yourname/envdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/envdiff.git
cd envdiff && go build -o envdiff .
```

---

## Usage

```bash
# Compare a base .env against one or more environment files
envdiff --base .env.example --compare .env.production

# Compare multiple environments at once
envdiff --base .env.example --compare .env.staging,.env.production

# Exit with a non-zero code on any diff (useful in CI pipelines)
envdiff --base .env.example --compare .env.production --strict
```

### Example Output

```
[MISSING]  DATABASE_URL   found in .env.example, missing in .env.production
[EXTRA]    DEBUG_MODE     found in .env.production, not in .env.example
[OK]       APP_PORT       matches across all files
```

---

## CI Integration

Add `envdiff` to your pipeline to catch configuration drift before deployment:

```yaml
# .github/workflows/envcheck.yml
- name: Check env files
  run: envdiff --base .env.example --compare .env.production --strict
```

The `--strict` flag ensures the job fails if any missing or extra keys are detected.

---

## Flags

| Flag        | Description                                      |
|-------------|--------------------------------------------------|
| `--base`    | Path to the reference `.env` file                |
| `--compare` | Comma-separated list of `.env` files to check    |
| `--strict`  | Exit with code 1 if any differences are found    |
| `--quiet`   | Suppress output, only return exit code           |

---

## License

MIT © [yourname](https://github.com/yourname)