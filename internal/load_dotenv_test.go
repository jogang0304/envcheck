package internal_test

import (
	"os"
	"testing"

	"github.com/jogang0304/envcheck/internal"
)

const envFileContent = `
# This is a comment
EXAMPLE_VAR=example_value

ANOTHER_VAR="another_value"
# Another comment
SOME_VAR='some_value'
last var =   last value
`

var expectedVars = map[string]string{
	"EXAMPLE_VAR": "example_value",
	"ANOTHER_VAR": "another_value",
	"SOME_VAR":    "some_value",
	"last var":    "last value",
}

func TestLoadDotenvFromFile(t *testing.T) {
	t.Run("Valid .env file", func(t *testing.T) {
		tempdir := t.TempDir()

		// Create a temporary .env file
		envFilePath := tempdir + "/.env"

		err := os.WriteFile(envFilePath, []byte(envFileContent), 0644)
		if err != nil {
			t.Fatalf("failed to create .env file: %v", err)
		}

		// Load the .env file
		err = internal.LoadDotenvFromFile(envFilePath)
		if err != nil {
			t.Fatalf("failed to load .env file: %v", err)
		}

		// Check if the environment variables are set correctly
		for key, expectedValue := range expectedVars {
			value, exists := os.LookupEnv(key)
			if !exists {
				t.Errorf("environment variable %s is not set", key)
				continue
			}
			if value != expectedValue {
				t.Errorf("environment variable %s has value %s, expected %s", key, value, expectedValue)
			}
		}
	})

	t.Run("Invalid .env file", func(t *testing.T) {
		invalidContents := []string{
			"INVALID_CONTENT",
			"VAR1==value1",
			"VAR2=value2=",
			"VAR3=value3=value4",
		}

		for _, content := range invalidContents {
			t.Run("Invalid content: "+content, func(t *testing.T) {

				tempdir := t.TempDir()

				// Create a temporary .env file with invalid content
				envFilePath := tempdir + "/.env"
				err := os.WriteFile(envFilePath, []byte(content), 0644)
				if err != nil {
					t.Fatalf("failed to create .env file: %v", err)
				}

				// Load the .env file
				err = internal.LoadDotenvFromFile(envFilePath)
				if err == nil {
					t.Fatalf("should have failed because of invalid .env file content")
				}
			})
		}
	})
}

func TestLoadDotenv(t *testing.T) {
	t.Run("Dotenv file exists in cwd", func(t *testing.T) {
		tempdir := t.TempDir()

		// Create a temporary .env file
		envFilePath := tempdir + "/.env"
		err := os.WriteFile(envFilePath, []byte(envFileContent), 0644)
		if err != nil {
			t.Fatalf("failed to create .env file: %v", err)
		}

		os.Chdir(tempdir)

		// Load the .env file
		err = internal.LoadDotenv()
		if err != nil {
			t.Fatalf("failed to load .env file: %v", err)
		}

		// Check if the environment variables are set correctly
		for key, expectedValue := range expectedVars {
			value, exists := os.LookupEnv(key)
			if !exists {
				t.Errorf("environment variable %s is not set", key)
				continue
			}
			if value != expectedValue {
				t.Errorf("environment variable %s has value %s, expected %s", key, value, expectedValue)
			}
		}
	})

	t.Run("Dotenv file does not exist in cwd", func(t *testing.T) {
		// Change to a temporary directory
		tempdir := t.TempDir()
		os.Chdir(tempdir)

		// Load the .env file
		err := internal.LoadDotenv()
		if err == nil {
			t.Fatalf("should have failed because no .env file exists")
		}
	})
}
