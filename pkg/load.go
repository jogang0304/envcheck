package envcheck

import (
	"errors"

	"github.com/jogang0304/envcheck/internal"
)

func Load() error {
	internal.LoadDotenv()
	c, err := internal.GetConfig()
	if err != nil {
		return errors.Join(err, errors.New("failed to get config"))
	}

	err = internal.PopulateUnsetVarsWithDefaults(c)
	if err != nil {
		return errors.Join(err, errors.New("failed to populate unset vars with defaults"))
	}

	err = internal.ValidateRequiredVars(c)
	if err != nil {
		return errors.Join(err, errors.New("failed to validate required vars"))
	}

	err = internal.ValidateVarTypes(c)
	if err != nil {
		return errors.Join(err, errors.New("failed to validate var types"))
	}

	err = internal.ValidateVarPatterns(c)
	if err != nil {
		return errors.Join(err, errors.New("failed to validate var patterns"))
	}

	return nil
}
