package internal_test

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/jogang0304/envcheck/internal"
)

func StringPtr(s string) *string {
	return &s
}

func testGetConfigWithError(t *testing.T, configFileContent *string, expectedErrorText string) {
	tempdir := t.TempDir()

	if configFileContent != nil {
		envFilePath := tempdir + "/.env.yaml"
		err := os.WriteFile(envFilePath, []byte(*configFileContent), 0o644)
		if err != nil {
			t.Fatalf("failed to create .env.config file: %v", err)
		}
	}

	err := os.Chdir(tempdir)
	if err != nil {
		t.Fatalf("failed to chdir to %s", tempdir)
	}

	_, err = internal.GetConfig()
	if err == nil {
		t.Fatal("expected to get an error")
	}

	hasExpectedError := strings.Contains(err.Error(), expectedErrorText)
	if !hasExpectedError {
		t.Fatalf("expected error \"%v\", got \"%v\"", expectedErrorText, err)
	}
}

func TestGetConfig(t *testing.T) {
	t.Run("Valid config file", func(t *testing.T) {
		const configFileContent = `
vars:
  - name: firstVar
    required: false
    type: int
    default_value: 0
  - name: secondVar
    required: true
    type: string
    pattern: .*secret.*
`
		expectedConfig := internal.Config{
			Vars: []internal.VarEntry{{
				Name:         "firstVar",
				Required:     false,
				Type:         internal.IntType,
				DefaultValue: 0,
				Pattern:      nil,
			}, {
				Name:         "secondVar",
				Required:     true,
				Type:         internal.StringType,
				Pattern:      StringPtr(".*secret.*"),
				DefaultValue: nil,
			}},
		}

		tempdir := t.TempDir()

		envFilePath := tempdir + "/.env.yaml"
		err := os.WriteFile(envFilePath, []byte(configFileContent), 0o644)
		if err != nil {
			t.Fatalf("failed to create .env.yaml file: %v", err)
		}

		err = os.Chdir(tempdir)
		if err != nil {
			t.Fatalf("failed to chdir to %s", tempdir)
		}

		c, err := internal.GetConfig()
		if err != nil {
			t.Fatalf("failed to read env config: %v", err)
		}

		configIsCorrect := reflect.DeepEqual(c, expectedConfig)

		if !configIsCorrect {
			t.Errorf("expected config: %v\n Actual config: %v", expectedConfig, c)
		}
	})

	t.Run("Invalid yaml file", func(t *testing.T) {
		t.Run("Incorrect yaml indentation", func(t *testing.T) {
			configFileContent := `
vars:
  - required: false
    name: firstVar
    type: int
    default_value: 0
  - name: secondVar
 required: true
`
			const expectedErrorText = "failed to unmarshal .env.config. Probably incorrect yaml structure"

			testGetConfigWithError(t, &configFileContent, expectedErrorText)
		})

		t.Run("No \"name\" field", func(t *testing.T) {
			configFileContent := `
vars:
  - required: false
    type: int
    default_value: 0
  - name: secondVar
`
			const expectedErrorText = "config has var without name"

			testGetConfigWithError(t, &configFileContent, expectedErrorText)
		})
	})

	t.Run("No .env.yaml file", func(t *testing.T) {
		const expectedErrorText = "failed to read .env.yaml"

		testGetConfigWithError(t, nil, expectedErrorText)
	})
}
