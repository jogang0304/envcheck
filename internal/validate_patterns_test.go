package internal_test

import (
	"os"
	"strings"
	"testing"

	"github.com/jogang0304/envcheck/internal"
)

func testValidatePatternsWithError(
	t *testing.T,
	config *internal.Config,
	expectedErrorText string,
) {
	err := internal.ValidatePatterns(config)
	if err == nil {
		t.Fatalf("expected to get an error")
	}

	if !strings.Contains(err.Error(), expectedErrorText) {
		t.Fatalf("expected to get \"%v\"\nin \"%v\"", expectedErrorText, err.Error())
	}
}

func TestValidatePatterns(t *testing.T) {
	t.Run("Correct env", func(t *testing.T) {
		keys := []string{"firstVar"}
		originalEnv := saveAndClearEnv(keys)
		defer restoreEnv(originalEnv)

		config := internal.Config{
			Vars: []internal.VarEntry{
				{
					Name:    "firstVar",
					Type:    "string",
					Pattern: StringPtr("[01]..$"),
				},
			},
		}

		if err := os.Setenv("firstVar", "123"); err != nil {
			t.Fatal("failed to set env firstVar")
		}

		err := internal.ValidatePatterns(&config)
		if err != nil {
			t.Fatalf("unexpected error \"%v\"", err)
		}
	})

	t.Run("Incorrect env", func(t *testing.T) {
		keys := []string{"firstVar"}
		originalEnv := saveAndClearEnv(keys)
		defer restoreEnv(originalEnv)

		config := internal.Config{
			Vars: []internal.VarEntry{
				{
					Name:    "firstVar",
					Type:    "string",
					Pattern: StringPtr("[01]..$"),
				},
			},
		}

		if err := os.Setenv("firstVar", "6581567231"); err != nil {
			t.Fatal("failed to set env firstVar")
		}

		const expectedErrorText = "variable firstVar does not match pattern [01]..$"

		testValidatePatternsWithError(t, &config, expectedErrorText)
	})

	t.Run("Incorrect pattern", func(t *testing.T) {
		keys := []string{"firstVar"}
		originalEnv := saveAndClearEnv(keys)
		defer restoreEnv(originalEnv)

		config := internal.Config{
			Vars: []internal.VarEntry{
				{
					Name:    "firstVar",
					Type:    "string",
					Pattern: StringPtr("{0,1}..$"),
				},
			},
		}

		if err := os.Setenv("firstVar", "i164187s"); err != nil {
			t.Fatal("failed to set env firstVar")
		}

		const expectedErrorText = "failed to compile regex for variable firstVar"

		testValidatePatternsWithError(t, &config, expectedErrorText)
	})

	t.Run("Incorrect var type", func(t *testing.T) {
		keys := []string{"firstVar"}
		originalEnv := saveAndClearEnv(keys)
		defer restoreEnv(originalEnv)

		config := internal.Config{
			Vars: []internal.VarEntry{
				{
					Name:    "firstVar",
					Type:    "int",
					Pattern: StringPtr("[01]..$"),
				},
			},
		}

		if err := os.Setenv("firstVar", "124"); err != nil {
			t.Fatal("failed to set env firstVar")
		}

		const expectedErrorText = "variable \"firstVar\" has type int, \"pattern\" is supported only for type string"

		testValidatePatternsWithError(t, &config, expectedErrorText)
	})
}
