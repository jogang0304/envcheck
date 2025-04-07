package internal

import (
	"fmt"
	"os"
	"strconv"
)

func ValidateVarTypes(c Config) error {
	// Validate var types
	vars := c.Vars

	for _, v := range vars {
		value, ok := os.LookupEnv(v.Name)
		if !ok {
			continue // unset vars do not have type. Wether they are allowed to be unset is handled in ValidateRequiredVars
		}

		if v.Type == StringType {
			continue // anything can be a string
		} else if v.Type == IntType {
			if _, err := strconv.Atoi(value); err != nil {
				return fmt.Errorf("var %s is not a valid int", v.Name)
			}
		} else if v.Type == FloatType {
			if _, err := strconv.ParseFloat(value, 64); err != nil {
				return fmt.Errorf("var %s is not a valid float", v.Name)
			}
		} else if v.Type == BoolType {
			if _, err := strconv.ParseBool(value); err != nil {
				return fmt.Errorf("var %s is not a valid bool", v.Name)
			}
		} else if v.Type == AnyType {
			continue
		} else {
			return fmt.Errorf("var %s has an unsupported type", v.Name)
		}
	}

	return nil
}
