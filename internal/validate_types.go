package internal

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

func ValidateTypes(c *Config) error {
	var typeError error = nil

	for _, v := range c.Vars {
		value, ok := os.LookupEnv(v.Name)
		if !ok {
			continue // unset vars do not have type. Wether they are allowed to be unset is handled in ValidateRequiredVars
		}

		if v.Type == StringType {
			continue // anything can be a string
		} else if v.Type == IntType {
			if _, err := strconv.Atoi(value); err != nil {
				typeError = errors.Join(fmt.Errorf("var %s is not a valid int", v.Name), typeError)
			}
		} else if v.Type == FloatType {
			if _, err := strconv.ParseFloat(value, 64); err != nil {
				typeError = errors.Join(fmt.Errorf("var %s is not a valid float", v.Name), typeError)
			}
		} else if v.Type == BoolType {
			if _, err := strconv.ParseBool(value); err != nil {
				typeError = errors.Join(fmt.Errorf("var %s is not a valid bool", v.Name), typeError)
			}
		} else if v.Type == AnyType {
			continue
		} else {
			typeError = errors.Join(fmt.Errorf("var %s has an unsupported type", v.Name), typeError)
		}
	}

	return typeError
}
