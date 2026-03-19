from fastmcp import FastMCP

from .client import PaprikaClient
from .config import get_settings, token_cache_path

mcp = FastMCP("Paprika Recipe Manager")

MEAL_TYPES = {0: "Breakfast", 1: "Lunch", 2: "Dinner", 3: "Snack"}


def _client() -> PaprikaClient:
    settings = get_settings()
    return PaprikaClient(
        email=settings.paprika_email,
        password=settings.paprika_password,
        token_cache_path=token_cache_path(),
    )


@mcp.tool()
def get_sync_status() -> dict:
    """Get sync status counters for all Paprika resource types.

    Returns change counters (not total counts) that increment on each
    modification. Useful for detecting which resource types have changed
    since the last sync.
    """
    return _client().get_sync_status()


@mcp.tool()
def list_recipes() -> list[dict]:
    """List all recipes as lightweight {uid, hash} pairs.

    Returns only uid and hash for each recipe — not full recipe data.
    Use get_recipe(uid) to fetch complete details for a specific recipe.
    The hash field can be used to detect changes without re-fetching unchanged recipes.
    """
    return _client().list_recipes()


@mcp.tool()
def get_recipe(uid: str) -> dict:
    """Get full details for a specific recipe by its UID.

    Returns the complete recipe object including name, ingredients, directions,
    notes, nutrition, timing, rating, categories, source, and metadata.

    Args:
        uid: The recipe's unique identifier (uppercase UUID4 format).
    """
    return _client().get_recipe(uid)


@mcp.tool()
def list_categories() -> list[dict]:
    """List all recipe categories.

    Returns category uid, name, order_flag, and parent_uid (for nested categories).
    Recipe objects reference categories by name, not UID.
    """
    return _client().list_categories()


@mcp.tool()
def list_grocery_lists() -> list[dict]:
    """List all grocery lists.

    Returns each list's uid, name, order_flag, is_default, and reminders_list fields.
    Use the uid from this response to filter grocery items by list.
    """
    return _client().list_grocery_lists()


@mcp.tool()
def list_grocery_items(list_uid: str, include_checked: bool = False) -> list[dict]:
    """List grocery items for a specific grocery list.

    Returns items with name, quantity, ingredient, aisle, purchased status,
    recipe reference, and list membership.

    Args:
        list_uid: Grocery list UID to filter items.
        include_checked: Include checked/purchased items when true.
    """
    items = _client().list_grocery_items()
    items = [item for item in items if item.get("list_uid") == list_uid]
    if not include_checked:
        items = [item for item in items if not item.get("purchased")]
    return items


@mcp.tool()
def list_meal_plans(start_date: str | None = None, end_date: str | None = None) -> list[dict]:
    """List meal plan entries, optionally filtered by date range.

    Returns meal entries with name, date, meal type, and optional recipe reference.
    Each entry includes a human-readable meal_type_name field in addition to the
    numeric type field (0=Breakfast, 1=Lunch, 2=Dinner, 3=Snack).

    Args:
        start_date: Optional start date (YYYY-MM-DD) inclusive.
        end_date: Optional end date (YYYY-MM-DD) inclusive.
    """
    meals = _client().list_meal_plans()
    if start_date or end_date:
        from datetime import date

        start = date.fromisoformat(start_date) if start_date else None
        end = date.fromisoformat(end_date) if end_date else None

        filtered: list[dict] = []
        for meal in meals:
            meal_date_str = meal.get("date", "").split(" ")[0]
            try:
                meal_date = date.fromisoformat(meal_date_str)
            except Exception:
                continue
            if start and meal_date < start:
                continue
            if end and meal_date > end:
                continue
            filtered.append(meal)
        meals = filtered
    for meal in meals:
        meal["meal_type_name"] = MEAL_TYPES.get(meal.get("type"), "Unknown")
    return meals


def main() -> None:
    settings = get_settings()
    mcp.run(transport="http", host="127.0.0.1", port=settings.paprika_port)


if __name__ == "__main__":
    main()
