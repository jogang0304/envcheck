package envcheck

import (
	"errors"

	"github.com/jogang0304/envcheck/internal"
)

func Load() error {
	err := internal.LoadDotenv()
	if err != nil {
		return errors.Join(errors.New("failed to load .env"), err)
	}

	c, err := internal.GetConfig()
	if err != nil {
		return errors.Join(errors.New("failed to get config"), err)
	}

	err = internal.PopulateUnsetVarsWithDefaults(&c)
	if err != nil {
		return errors.Join(errors.New("failed to populate unset vars with defaults"), err)
	}

	err = internal.ValidateRequiredVars(&c)
	if err != nil {
		return errors.Join(errors.New("failed to validate required vars"), err)
	}

	err = internal.ValidateVarTypes(&c)
	if err != nil {
		return errors.Join(errors.New("failed to validate var types"), err)
	}

	err = internal.ValidateVarPatterns(&c)
	if err != nil {
		return errors.Join(errors.New("failed to validate var patterns"), err)
	}

	return nil
}
