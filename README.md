# Paprika MCP Server

FastMCP server for Paprika Recipe Manager.

## Quick Install (macOS)

```bash
curl -sSL https://raw.githubusercontent.com/aarons22/paprika-mcp/main/install.sh | bash
```

Then run:

```bash
paprika-mcp setup
paprika-mcp install
```

If `paprika-mcp` isn’t on your PATH, use:

```bash
$HOME/.local/bin/paprika-mcp --help
```

## CLI Commands

- `paprika-mcp setup`
- `paprika-mcp run`
- `paprika-mcp install`
- `paprika-mcp uninstall`
- `paprika-mcp status`
- `paprika-mcp logs`

## Config

- `~/Library/Application Support/paprika-mcp/config.toml`
- Token cache: `~/Library/Application Support/paprika-mcp/.paprika_token.json`

## Logs

- `~/Library/Logs/paprika-mcp.out.log`
- `~/Library/Logs/paprika-mcp.err.log`

## Notes

- All operations are **read-only**
- The Paprika API is unofficial and undocumented; see `API_REFERENCE.md` for details
- Tokens are cached on disk and refreshed automatically on expiry (401)
