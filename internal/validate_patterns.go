package internal

import (
	"fmt"
	"os"
	"regexp"
)

/*
If var.Pattern is not empty, it is a regex with which var value should be checked.
If var type is not a string, it should be converted to string and then checked with regex.
*/
func ValidateVarPatterns(c Config) error {
	// Validate the patterns of the variables
	vars := c.Vars

	for _, v := range vars {
		if v.Pattern == "" {
			continue
		}

		value, ok := os.LookupEnv(v.Name)
		if !ok {
			// If the variable is not set, skip the validation
			continue
		}

		if v.Type != StringType {
			// Convert the value to string
			value = fmt.Sprintf("%v", value)
		}

		// Check if the value matches the pattern
		matched, err := regexp.MatchString(v.Pattern, value)
		if err != nil {
			return fmt.Errorf("failed to compile regex for variable %s: %w", v.Name, err)
		}
		if !matched {
			return fmt.Errorf("variable %s does not match the pattern %s", v.Name, v.Pattern)
		}
	}
	return nil
}
