package client

import (
	"encoding/json"
	"fmt"
)

// PaginationType identifies the pagination style for auto-pagination.
type PaginationType string

const (
	PaginationPageBased   PaginationType = "page"
	PaginationOffsetBased PaginationType = "offset"
	PaginationCursorBased PaginationType = "cursor"

	// maxPaginationPages is the safety limit preventing infinite loops.
	maxPaginationPages = 100
)

// PaginationConfig holds parameters for auto-paginating a request.
type PaginationConfig struct {
	Type        PaginationType
	PageParam   string // page number or offset param name
	SizeParam   string // per_page / limit param name
	CursorParam string // cursor / after / before param name
	PageSize    int    // items per page (0 = use API default)
}

// FetchAll auto-paginates through all pages of a paginated API endpoint and
// returns the aggregated raw JSON results. It respects a safety limit of
// maxPaginationPages to prevent infinite loops.
//
// client is used to perform each HTTP request. method, path, pathParams and
// baseQuery are passed through to client.Do on every request; pagination
// parameters are appended/updated automatically.
func FetchAll(c *Client, method, path string, pathParams map[string]string, baseQuery map[string]string, cfg PaginationConfig) ([]byte, error) {
	var allItems []interface{}
	query := make(map[string]string, len(baseQuery))
	for k, v := range baseQuery {
		query[k] = v
	}

	switch cfg.Type {
	case PaginationPageBased:
		return fetchAllPageBased(c, method, path, pathParams, query, cfg, &allItems)
	case PaginationOffsetBased:
		return fetchAllOffsetBased(c, method, path, pathParams, query, cfg, &allItems)
	case PaginationCursorBased:
		return fetchAllCursorBased(c, method, path, pathParams, query, cfg, &allItems)
	default:
		return nil, fmt.Errorf("unknown pagination type: %s", cfg.Type)
	}
}

func fetchAllPageBased(c *Client, method, path string, pathParams map[string]string, query map[string]string, cfg PaginationConfig, allItems *[]interface{}) ([]byte, error) {
	startPage := 1
	if cfg.PageParam != "" {
		if v, ok := query[cfg.PageParam]; ok {
			if n, scanErr := fmt.Sscanf(v, "%d", &startPage); n != 1 || scanErr != nil {
				startPage = 1
			}
		}
	}
	page := startPage
	for i := 0; i < maxPaginationPages; i++ {
		if cfg.PageParam != "" {
			query[cfg.PageParam] = fmt.Sprintf("%d", page)
		}
		if cfg.SizeParam != "" && cfg.PageSize > 0 {
			query[cfg.SizeParam] = fmt.Sprintf("%d", cfg.PageSize)
		}

		resp, err := c.Do(method, path, pathParams, query, nil)
		if err != nil {
			return nil, err
		}

		items, done, err := extractItems(resp)
		if err != nil || done {
			break
		}
		*allItems = append(*allItems, items...)
		page++
	}
	return json.Marshal(*allItems)
}

func fetchAllOffsetBased(c *Client, method, path string, pathParams map[string]string, query map[string]string, cfg PaginationConfig, allItems *[]interface{}) ([]byte, error) {
	offset := 0
	limit := cfg.PageSize
	if limit <= 0 && cfg.SizeParam != "" {
		if v, ok := query[cfg.SizeParam]; ok {
			if n, scanErr := fmt.Sscanf(v, "%d", &limit); n != 1 || scanErr != nil {
				limit = 0
			}
		}
	}
	if limit <= 0 {
		limit = 50
	}
	for i := 0; i < maxPaginationPages; i++ {
		if cfg.PageParam != "" {
			query[cfg.PageParam] = fmt.Sprintf("%d", offset)
		}
		if cfg.SizeParam != "" {
			query[cfg.SizeParam] = fmt.Sprintf("%d", limit)
		}

		resp, err := c.Do(method, path, pathParams, query, nil)
		if err != nil {
			return nil, err
		}

		items, done, err := extractItems(resp)
		if err != nil || done {
			break
		}
		*allItems = append(*allItems, items...)
		offset += limit
	}
	return json.Marshal(*allItems)
}

func fetchAllCursorBased(c *Client, method, path string, pathParams map[string]string, query map[string]string, cfg PaginationConfig, allItems *[]interface{}) ([]byte, error) {
	cursor := ""
	for i := 0; i < maxPaginationPages; i++ {
		if cursor != "" && cfg.CursorParam != "" {
			query[cfg.CursorParam] = cursor
		}

		resp, err := c.Do(method, path, pathParams, query, nil)
		if err != nil {
			return nil, err
		}

		items, done, err := extractItems(resp)
		if err != nil || done {
			break
		}
		*allItems = append(*allItems, items...)

		// Extract next cursor using canonical field conventions.
		cursor = extractNextCursor(resp)
		if cursor == "" {
			break
		}
	}
	return json.Marshal(*allItems)
}

// extractItems parses a JSON response and returns the items array.
// Returns done=true when the response is empty or contains no items.
func extractItems(data []byte) ([]interface{}, bool, error) {
	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, true, err
	}
	switch v := raw.(type) {
	case []interface{}:
		if len(v) == 0 {
			return nil, true, nil
		}
		return v, false, nil
	case map[string]interface{}:
		// Look for a "data", "items", "results", or "records" key.
		for _, key := range []string{"data", "items", "results", "records"} {
			if arr, ok := v[key].([]interface{}); ok {
				if len(arr) == 0 {
					return nil, true, nil
				}
				return arr, false, nil
			}
		}
		// No recognised envelope key found — treat as terminal (e.g. error response).
		return nil, true, nil
	default:
		return nil, true, nil
	}
}

// extractNextCursor tries canonical cursor field names from a JSON response:
// next_cursor, nextCursor, after, meta.next.
func extractNextCursor(data []byte) string {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return ""
	}

	// Direct fields: next_cursor, nextCursor, after.
	for _, field := range []string{"next_cursor", "nextCursor", "after"} {
		if v, ok := raw[field].(string); ok && v != "" {
			return v
		}
	}

	// meta.next
	if meta, ok := raw["meta"].(map[string]interface{}); ok {
		if v, ok := meta["next"].(string); ok && v != "" {
			return v
		}
	}

	return ""
}
