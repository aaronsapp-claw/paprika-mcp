# paprika-mcp

MCP server for Paprika Recipe Manager, built with [FastMCP](https://github.com/jlowin/fastmcp).

## Setup

1. Install dependencies:
   ```bash
   pip install -r requirements.txt
   ```

2. Copy `.env.example` to `.env` and fill in your Paprika credentials:
   ```bash
   cp .env.example .env
   ```

## Claude Desktop Configuration

Add to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "paprika": {
      "command": "python",
      "args": ["/path/to/paprika-mcp/server.py"]
    }
  }
}
```

## Available Tools

| Tool | Description |
|------|-------------|
| `get_sync_status` | Get change counters for all resource types |
| `list_recipes` | List all recipes as lightweight `{uid, hash}` pairs |
| `get_recipe(uid)` | Get full recipe details by UID |
| `list_categories` | List all recipe categories |
| `list_grocery_lists` | List all grocery lists |
| `list_grocery_items(list_uid?)` | List grocery items, optionally filtered by list |
| `list_meal_plans` | List all meal plan entries |

## Notes

- All operations are **read-only**
- The Paprika API is unofficial and undocumented; see `API_REFERENCE.md` for details
- Tokens are cached in memory and refreshed automatically on expiry (401)
