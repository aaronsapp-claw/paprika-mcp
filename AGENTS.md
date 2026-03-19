# AGENTS

## End-to-End Validation (Required)
When making changes that affect Paprika auth, network handling, or the MCP server runtime, run a local end-to-end validation before declaring success.

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
