package internal_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/jogang0304/envcheck/internal"
)

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

func TestPopulateUnsetRequiredVarsWithDefaults(t *testing.T) {
	// Define a test case
	testCases := []struct {
		name           string
		config         internal.Config
		presetVars     map[string]any
		expectedResult map[string]any
	}{
		{
			name: "All vars have defaults",
			config: internal.Config{
				Vars: []internal.VarEntry{
					{
						Name:         "EXAMPLE_VAR",
						Required:     true,
						Type:         internal.StringType,
						DefaultValue: "default_example",
					},
					{
						Name:         "ANOTHER_VAR",
						Required:     false,
						Type:         internal.StringType,
						DefaultValue: "default_another",
					},
				},
			},
			presetVars: map[string]any{
				"EXAMPLE_VAR": "preset_value",
			},
			expectedResult: map[string]any{
				"EXAMPLE_VAR": "preset_value",
				"ANOTHER_VAR": "default_another",
			},
		},
		{
			name: "No vars have defaults",
			config: internal.Config{
				Vars: []internal.VarEntry{
					{Name: "EXAMPLE_VAR", Required: true, Type: internal.StringType},
					{Name: "ANOTHER_VAR", Required: false, Type: internal.StringType},
				},
			},
			presetVars: map[string]any{
				"EXAMPLE_VAR": "preset_value",
			},
			expectedResult: map[string]any{
				"EXAMPLE_VAR": "preset_value",
				"ANOTHER_VAR": nil,
			},
		},
		{
			name: "Some vars have defaults",
			config: internal.Config{
				Vars: []internal.VarEntry{
					{
						Name:         "EXAMPLE_VAR",
						Required:     true,
						Type:         internal.StringType,
						DefaultValue: "default_example",
					},
					{Name: "ANOTHER_VAR", Required: false, Type: internal.StringType},
					{
						Name:         "SOME_VAR",
						Required:     true,
						Type:         internal.StringType,
						DefaultValue: "default_some",
					},
				},
			},
			presetVars: map[string]any{
				"EXAMPLE_VAR": "preset_value",
			},
			expectedResult: map[string]any{
				"EXAMPLE_VAR": "preset_value",
				"ANOTHER_VAR": nil,
				"SOME_VAR":    "default_some",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var keys []string
			for _, v := range tc.config.Vars {
				keys = append(keys, v.Name)
			}
			// Save and clear the environment variables
			originalEnv := saveAndClearEnv(keys)
			defer restoreEnv(originalEnv)

			// Set preset environment variables
			for key, value := range tc.presetVars {
				t.Setenv(key, fmt.Sprintf("%v", value))
			}

			err := internal.PopulateUnsetVarsWithDefaults(&tc.config)
			if err != nil {
				t.Fatalf("failed to populate unset vars with defaults: %v", err)
			}

			// Check if the environment variables are set correctly
			for v := range tc.expectedResult {
				value, exists := os.LookupEnv(v)
				if !exists && tc.expectedResult[v] != nil {
					t.Errorf("environment variable %s is not set", v)
					continue
				}
				if tc.expectedResult[v] == nil && exists {
					t.Errorf("environment variable %s must not be set. But it is %s", v, value)
					continue
				}
				if tc.expectedResult[v] == nil && !exists {
					continue
				}
				expectedValue := fmt.Sprintf("%v", tc.expectedResult[v])
				if value != expectedValue {
					t.Errorf(
						"environment variable %s has value %s, expected %s",
						v,
						value,
						expectedValue,
					)
				}
			}
		})
	}
}
