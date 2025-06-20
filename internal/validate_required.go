package internal

import (
	"errors"
	"fmt"
	"os"
)

// Validate required vars
func ValidateRequired(c *Config) error {
	var requireError error = nil

	for _, v := range c.Vars {
		if v.Required {
			if _, ok := os.LookupEnv(v.Name); !ok {
				requireError = errors.Join(
					fmt.Errorf("required var %s is not set", v.Name),
					requireError,
				)
			}
		}
	}

	return requireError
}
