package internal

import (
	"errors"
	"fmt"
	"os"
)

/*
This function iterates through c["vars"] and for each var that is not set, it populates the var with its default value. "Populate" means to set os.env[var.Name].
*/
func PopulateUnsetVarsWithDefaults(c Config) error {
	vars := c.Vars
	for _, v := range vars {
		_, isSet := os.LookupEnv(v.Name)
		if !isSet {
			if v.DefaultValue != nil {
				err := os.Setenv(v.Name, fmt.Sprintf("%v", v.DefaultValue))
				if err != nil {
					return errors.Join(errors.New("Failed to set env var"), err)
				}
			}
		}
	}
	return nil
}
