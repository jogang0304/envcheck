package internal

import (
	"fmt"
	"os"
)

func ValidateRequiredVars(c Config) error {
	// Validate required vars
	vars := c.Vars

	for _, v := range vars {
		if v.Required {
			if _, ok := os.LookupEnv(v.Name); !ok {
				return fmt.Errorf("required var %s is not set", v.Name)
			}
		}
	}

	return nil
}
