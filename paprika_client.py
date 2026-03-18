import base64
import gzip
import json
import os
from typing import Optional

import httpx
from dotenv import load_dotenv

load_dotenv()

BASE_URL = "https://www.paprikaapp.com/api"


class PaprikaClient:
    """HTTP client for the Paprika Recipe Manager API."""

    def __init__(self) -> None:
        self.email = os.environ["PAPRIKA_EMAIL"]
        self.password = os.environ["PAPRIKA_PASSWORD"]
        self._token: Optional[str] = None

    def _authenticate(self) -> str:
        """Obtain a bearer token using V1 Basic Auth + form data login."""
        credentials = base64.b64encode(
            f"{self.email}:{self.password}".encode()
        ).decode()
        response = httpx.post(
            f"{BASE_URL}/v1/account/login/",
            headers={
                "Authorization": f"Basic {credentials}",
                "Content-Type": "application/x-www-form-urlencoded",
            },
            data={"email": self.email, "password": self.password},
            timeout=30,
        )
        response.raise_for_status()
        self._token = response.json()["result"]["token"]
        return self._token

    def _get_token(self) -> str:
        if not self._token:
            self._authenticate()
        return self._token  # type: ignore[return-value]

    def _request(self, method: str, path: str, **kwargs) -> dict:
        """Make an authenticated request, retrying once on 401."""
        token = self._get_token()
        response = httpx.request(
            method,
            f"{BASE_URL}{path}",
            headers={"Authorization": f"Bearer {token}"},
            timeout=30,
            **kwargs,
        )
        if response.status_code == 401:
            self._token = None
            token = self._authenticate()
            response = httpx.request(
                method,
                f"{BASE_URL}{path}",
                headers={"Authorization": f"Bearer {token}"},
                timeout=30,
                **kwargs,
            )
        response.raise_for_status()

        # Some responses are gzip-compressed; check magic bytes
        content = response.content
        if content[:2] == b"\x1f\x8b":
            content = gzip.decompress(content)
        return json.loads(content)

    # --- Public API methods ---

    def get_sync_status(self) -> dict:
        """Return change counters for all Paprika resource types."""
        return self._request("GET", "/v2/sync/status/")["result"]

    def list_recipes(self) -> list:
        """Return lightweight list of {uid, hash} pairs for all recipes."""
        return self._request("GET", "/v2/sync/recipes/")["result"]

    def get_recipe(self, uid: str) -> dict:
        """Return full details for a single recipe by UID."""
        return self._request("GET", f"/v2/sync/recipe/{uid}/")["result"]

    def list_categories(self) -> list:
        """Return all recipe categories."""
        return self._request("GET", "/v2/sync/categories/")["result"]

    def list_grocery_lists(self) -> list:
        """Return all grocery lists."""
        return self._request("GET", "/v2/sync/grocerylists/")["result"]

    def list_grocery_items(self) -> list:
        """Return all grocery items across all lists."""
        return self._request("GET", "/v2/sync/groceries/")["result"]

    def list_meal_plans(self) -> list:
        """Return all meal plan entries."""
        return self._request("GET", "/v2/sync/meals/")["result"]
