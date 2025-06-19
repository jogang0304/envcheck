package envcheck_test

import (
	"os"
	"strings"
	"testing"

	envcheck "github.com/jogang0304/envcheck/pkg"
)

func checkEnvVars(t *testing.T, expectedVars map[string]string) {
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
}

func createFilesInTempDir(t *testing.T, configFileContent, envFileContent string) {
	tempdir := t.TempDir()

	envYamlFilePath := tempdir + "/.env.yaml"
	err := os.WriteFile(envYamlFilePath, []byte(configFileContent), 0644)
	if err != nil {
		t.Fatalf("failed to create .env.yaml file: %v", err)
	}

	envFilePath := tempdir + "/.env"
	err = os.WriteFile(envFilePath, []byte(envFileContent), 0644)
	if err != nil {
		t.Fatalf("failed to create .env file: %v", err)
	}

	os.Chdir(tempdir)
}

func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

func saveAndClearEnv(keys []string) map[string]*string {
	original := make(map[string]*string)
	for _, key := range keys {
		if val, ok := os.LookupEnv(key); ok {
			v := val // копия строки
			original[key] = &v
			_ = os.Unsetenv(key)
		} else {
			original[key] = nil
		}
	}
	return original
}

func restoreEnv(env map[string]*string) {
	for key, val := range env {
		if val == nil {
			_ = os.Unsetenv(key)
		} else {
			_ = os.Setenv(key, *val)
		}
	}
}

func testLoadWithoutErrors(t *testing.T, configFileContent, envFileContent string, expectedVars map[string]string) {
	originalEnv := saveAndClearEnv(Keys(expectedVars))
	defer restoreEnv(originalEnv)

	createFilesInTempDir(t, configFileContent, envFileContent)

	err := envcheck.Load()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	checkEnvVars(t, expectedVars)
}

func testLoadWithErrors(t *testing.T, configFileContent, envFileContent, expectedError string, varsToClear []string) {
	originalEnv := saveAndClearEnv(varsToClear)
	defer restoreEnv(originalEnv)

	createFilesInTempDir(t, configFileContent, envFileContent)

	err := envcheck.Load()
	if err == nil {
		t.Fatalf("err == nil, but expected to fail with error \"%v\"", expectedError)
	}

	if !strings.Contains(err.Error(), expectedError) {
		t.Fatalf("expected \"%v\"\nto contain \"%v\"", err.Error(), expectedError)
	}
}

func TestLoad(t *testing.T) {
	t.Run("Correct .env and correct .env.yaml", func(t *testing.T) {
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
		const envFileContent = `
thirdVar=test123
fourth_var = test4
fifth VAR=  test 5
secondVar=123secret456
`

		var expectedVars = map[string]string{
			"thirdVar":   "test123",
			"fourth_var": "test4",
			"fifth VAR":  "test 5",
			"firstVar":   "0",
			"secondVar":  "123secret456",
		}

		testLoadWithoutErrors(t, configFileContent, envFileContent, expectedVars)
	})

	t.Run("Required vars are not set", func(t *testing.T) {
		t.Run(".env.yaml has default value", func(t *testing.T) {
			const configFileContent = `
vars:
  - name: secondVar
    required: true
    type: string
    default_value: def123
`
			const envFileContent = `
firstVar=test
`

			var expectedVars = map[string]string{
				"firstVar":  "test",
				"secondVar": "def123",
			}

			testLoadWithoutErrors(t, configFileContent, envFileContent, expectedVars)
		})

		t.Run(".env.yaml does not have default value", func(t *testing.T) {
			const configFileContent = `
vars:
  - name: secondVar
    required: true
    type: string
`
			const envFileContent = `
firstVar=test
`
			const expectedError = "failed to validate required vars"

			var varsToClean = []string{
				"firstVar", "secondVar",
			}

			testLoadWithErrors(t, configFileContent, envFileContent, expectedError, varsToClean)
		})
	})

	t.Run("Var has incorrect type", func(t *testing.T) {
		const configFileContent = `
vars:
  - name: firstVar
    required: true
    type: int
`
		const envFileContent = `
firstVar=test
`
		const expectedError = "failed to validate var types"

		var varsToClean = []string{
			"firstVar",
		}

		testLoadWithErrors(t, configFileContent, envFileContent, expectedError, varsToClean)
	})

	t.Run("Var has incorrect pattern", func(t *testing.T) {
		const configFileContent = `
vars:
  - name: firstVar
    required: true
    type: string
    pattern: .*123..$
`
		const envFileContent = `
firstVar=123
`
		const expectedError = "failed to validate var patterns"

		var varsToClean = []string{
			"firstVar",
		}

		testLoadWithErrors(t, configFileContent, envFileContent, expectedError, varsToClean)
	})

	t.Run(".env.yaml is incorrect", func(t *testing.T) {
		const configFileContent = `
vars:
  - required: true
    type: string
    pattern: .*123..$
`
		const envFileContent = `
firstVar=123
`
		const expectedError = "failed to get config"

		var varsToClean = []string{
			"firstVar",
		}

		testLoadWithErrors(t, configFileContent, envFileContent, expectedError, varsToClean)
	})
}
