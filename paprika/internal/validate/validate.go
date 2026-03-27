package validate

import (
	"fmt"
	"strings"
)

// Enum checks that val is one of the allowed values.
// Returns nil if val is empty (flag not set). Call only after cmd.Flags().Changed().
func Enum(flagName, val string, allowed []string) error {
	for _, a := range allowed {
		if a == val {
			return nil
		}
	}
	return fmt.Errorf("invalid value %q for --%s: must be one of: %s", val, flagName, strings.Join(allowed, ", "))
}
