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

	config, err := internal.GetConfig()
	if err != nil {
		return errors.Join(errors.New("failed to get config"), err)
	}

	err = internal.PopulateUnsetVarsWithDefaults(&config)
	if err != nil {
		return errors.Join(errors.New("failed to populate unset vars with defaults"), err)
	}

	err = internal.ValidateRequired(&config)
	if err != nil {
		return errors.Join(errors.New("failed to validate required vars"), err)
	}

	err = internal.ValidateTypes(&config)
	if err != nil {
		return errors.Join(errors.New("failed to validate var types"), err)
	}

	err = internal.ValidatePatterns(&config)
	if err != nil {
		return errors.Join(errors.New("failed to validate var patterns"), err)
	}

	return nil
}
