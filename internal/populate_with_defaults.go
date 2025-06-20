package internal

import (
	"errors"
	"fmt"
	"os"
)

/*
This function iterates through c["vars"] and for each var that is not set, it populates the var with its default value. "Populate" means to set os.env[var.Name].
*/
func PopulateUnsetVarsWithDefaults(c *Config) error {
	var populationError error = nil

	for _, v := range c.Vars {
		_, isSet := os.LookupEnv(v.Name)
		if !isSet {
			if v.DefaultValue != nil {
				err := os.Setenv(v.Name, fmt.Sprintf("%v", v.DefaultValue))
				if err != nil {
					populationError = errors.Join(
						errors.Join(errors.New("failed to set env var"), err),
						populationError,
					)
				}
			}
		}
	}
	return populationError
}
