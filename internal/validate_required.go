package internal

import (
	"fmt"
	"os"
)

// Validate required vars
func ValidateRequiredVars(c *Config) error {
	for _, v := range c.Vars {
		if v.Required {
			if _, ok := os.LookupEnv(v.Name); !ok {
				return fmt.Errorf("required var %s is not set", v.Name)
			}
		}
	}

	return nil
}
