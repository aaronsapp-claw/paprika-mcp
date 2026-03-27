package output

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/olekukonko/tablewriter"
)

// Print writes data to stdout. When jsonMode is true the raw JSON is printed
// compactly; otherwise it is pretty-printed with indentation.
func Print(data interface{}, jsonMode bool) error {
	var encoded []byte
	var err error

	if jsonMode {
		encoded, err = json.Marshal(data)
	} else {
		encoded, err = json.MarshalIndent(data, "", "  ")
	}
	if err != nil {
		return fmt.Errorf("marshalling output: %w", err)
	}

	_, err = fmt.Fprintln(os.Stdout, string(encoded))
	return err
}

// PrintTable renders JSON data as a human-readable table.
// Arrays of objects are rendered as columnar tables with one row per object.
// Single objects are rendered as two-column key-value tables.
// Nested objects within rows and arrays of non-objects fall back to
// pretty-printed JSON output.
func PrintTable(data []byte, noColor bool) error {
	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		// Fallback: print raw
		_, err = fmt.Fprintln(os.Stdout, string(data))
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	if noColor {
		table.SetBorder(true)
	}

	switch v := raw.(type) {
	case []interface{}:
		if len(v) == 0 {
			fmt.Fprintln(os.Stdout, "(no results)")
			return nil
		}
		// Collect headers from the first element.
		first, ok := v[0].(map[string]interface{})
		if !ok {
			// Not an array of objects; fall back to JSON.
			encoded, _ := json.MarshalIndent(raw, "", "  ")
			fmt.Fprintln(os.Stdout, string(encoded))
			return nil
		}
		headers := make([]string, 0, len(first))
		for k := range first {
			headers = append(headers, k)
		}
		sort.Strings(headers)
		table.SetHeader(headers)
		for _, item := range v {
			row, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			cols := make([]string, len(headers))
			for i, h := range headers {
				cols[i] = fmt.Sprintf("%v", row[h])
			}
			table.Append(cols)
		}
		table.Render()

	case map[string]interface{}:
		table.SetHeader([]string{"Key", "Value"})
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			table.Append([]string{k, fmt.Sprintf("%v", v[k])})
		}
		table.Render()

	default:
		fmt.Fprintln(os.Stdout, string(data))
	}

	return nil
}
