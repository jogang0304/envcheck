package internal

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type SupportedVarType string

const (
	StringType SupportedVarType = "string"
	IntType    SupportedVarType = "int"
	FloatType  SupportedVarType = "float"
	BoolType   SupportedVarType = "bool"
	AnyType    SupportedVarType = "any"
)

type VarEntry struct {
	Name         string           `yaml:"name"`
	Required     bool             `yaml:"required"`
	Type         SupportedVarType `yaml:"type"`
	DefaultValue any              `yaml:"default_value"`
	Pattern      *string          `yaml:"pattern"`
}

type Config struct {
	Vars []VarEntry `yaml:"vars"`
}

func checkRequiredFields(config *Config) error {
	for _, v := range config.Vars {
		if v.Name == "" {
			return errors.New("config has var without name")
		}
	}

	return nil
}

/*
This function reads .env.yaml file and returns config.
*/
func GetConfig() (Config, error) {
	data, err := os.ReadFile(".env.yaml")
	if err != nil {
		return Config{}, errors.Join(errors.New("failed to read .env.yaml"), err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Config{}, errors.Join(
			errors.New("failed to unmarshal .env.yaml. Probably incorrect yaml structure"),
			err,
		)
	}

	err = checkRequiredFields(&config)
	if err != nil {
		return Config{}, errors.Join(errors.New("failed to check required fields"), err)
	}

	return config, nil
}
