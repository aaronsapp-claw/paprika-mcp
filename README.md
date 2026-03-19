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

## Available Tools

| Tool | Description |
|------|-------------|
| `get_sync_status` | Get change counters for all resource types |
| `list_recipes` | List all recipes as lightweight `{uid, hash}` pairs |
| `get_recipe(uid)` | Get full recipe details by UID |
| `list_categories` | List all recipe categories |
| `list_grocery_lists` | List all grocery lists |
| `list_grocery_items(list_uid, include_checked?)` | List grocery items for a specific list |
| `list_meal_plans(start_date?, end_date?)` | List meal plan entries, optionally filtered by date |

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
