package internal

import (
	"fmt"
	"os"
	"regexp"
)

/*
If var.Pattern is not empty, it is a regex against which var value should be checked.
*/
func ValidateVarPatterns(c *Config) error {
	for _, v := range c.Vars {
		if v.Pattern == nil {
			continue
		}

		value, ok := os.LookupEnv(v.Name)
		if !ok {
			// If the variable is not set, skip the validation
			continue
		}

		if v.Type != StringType {
			return fmt.Errorf(
				"variable \"%s\" has type %s, \"pattern\" is supported only for type %s",
				v.Name,
				v.Type,
				StringType,
			)
		}

		// Check if the value matches the pattern
		matched, err := regexp.MatchString(*v.Pattern, value)
		if err != nil {
			return fmt.Errorf("failed to compile regex for variable %s: %w", v.Name, err)
		}
		if !matched {
			return fmt.Errorf("variable %s does not match pattern %v", v.Name, *v.Pattern)
		}
	}
	return nil
}
