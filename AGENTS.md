# AGENTS

## End-to-End Validation (Required)
When making changes that affect Paprika auth, network handling, or the MCP server runtime, run a local end-to-end validation before declaring success.

## Documentation Update (Required)
When tool signatures or behavior change, update the README and tool documentation (including `API_REFERENCE.md`) in the same change.

**Minimum E2E sequence:**
1. `python - <<'PY'`
   from paprika_mcp.client import PaprikaClient
   from paprika_mcp.config import get_settings, token_cache_path

   settings = get_settings()
   client = PaprikaClient(
       email=settings.paprika_email,
       password=settings.paprika_password,
       token_cache_path=token_cache_path(),
   )
   status = client.get_sync_status()
   print("status keys:", list(status.keys()) if isinstance(status, dict) else "n/a")
   PY

2. `paprika-mcp run` and verify the server starts without errors in logs.

**Notes:**
- If this is your first local run, execute `paprika-mcp setup` to configure credentials.

## CLI Code Generation (paprika/)

The Go CLI in `paprika/` is generated from `openapi.yaml` using [CommandSpec](https://github.com/theaiteam-dev/commandspec):

```bash
go install github.com/theaiteam-dev/commandspec@latest
commandspec init --schema ./openapi.yaml --name paprika --output-dir ./paprika
```

**Regenerating wipes all custom code unless it is wrapped in preservation markers:**

```go
// commandspec:custom:start
... your custom code ...
// commandspec:custom:end
```

Use `commandspec update` (not `init`) when the spec changes after the initial generation — it respects these markers. `init` always generates fresh files with no preservation.

**Files with custom markers (must not be removed):**

| File | What's custom |
|------|---------------|
| `paprika/cmd/account-login_login.go` | Entire file — `--email`/`--password` flags, hardcoded auth URL, auto-save token logic |
| `paprika/internal/config.go` | `Save()` function — writes token to `~/.config/paprika/config.yaml` |

After regenerating, also fix:
- `go.mod` module name → `github.com/aarons22/paprika-tools/paprika`
- All `import "paprika/..."` → `import "github.com/aarons22/paprika-tools/paprika/..."`
- Remove unused `"encoding/base64"` import from `internal/client/client.go`
- Patch `account-login.go` group name from `account-login` → `account`
