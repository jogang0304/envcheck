package internal_test

import (
	"os"
	"strings"
	"testing"

	"reflect"

	"github.com/jogang0304/envcheck/internal"
)

func StringPtr(s string) *string {
	return &s
}

func testGetConfigWithError(t *testing.T, configFileContent *string, expectedErrorText string) {
	tempdir := t.TempDir()

	if configFileContent != nil {
		envFilePath := tempdir + "/.env.yaml"
		err := os.WriteFile(envFilePath, []byte(*configFileContent), 0644)
		if err != nil {
			t.Fatalf("failed to create .env.config file: %v", err)
		}
	}

	os.Chdir(tempdir)

	_, err := internal.GetConfig()
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
		var expectedConfig internal.Config = internal.Config{
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
		err := os.WriteFile(envFilePath, []byte(configFileContent), 0644)
		if err != nil {
			t.Fatalf("failed to create .env.yaml file: %v", err)
		}

		os.Chdir(tempdir)

		c, err := internal.GetConfig()
		if err != nil {
			t.Fatalf("failed to read env config: %v", err)
		}

		configIsCorrect := reflect.DeepEqual(c, expectedConfig)

		if !configIsCorrect {
			t.Errorf("Expected config: %v\n Actual config: %v", expectedConfig, c)
		}
	})

	t.Run("Invalid yaml file", func(t *testing.T) {
		t.Run("Incorrect yaml indentation", func(t *testing.T) {
			var configFileContent = `
vars:
  - required: false
    name: firstVar
    type: int
    default_value: 0
  - name: secondVar
 required: true
`
			const expectedErrorText = "Failed to unmarshal .env.config. Probably incorrect yaml structure"

			testGetConfigWithError(t, &configFileContent, expectedErrorText)
		})

		t.Run("No \"name\" field", func(t *testing.T) {
			var configFileContent = `
vars:
  - required: false
    type: int
    default_value: 0
  - name: secondVar
`
			const expectedErrorText = "Config has var without name"

			testGetConfigWithError(t, &configFileContent, expectedErrorText)
		})
	})

	t.Run("No .env.yaml file", func(t *testing.T) {
		const expectedErrorText = "Failed to read .env.yaml"

		testGetConfigWithError(t, nil, expectedErrorText)
	})
}
