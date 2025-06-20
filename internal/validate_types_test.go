package internal_test

import (
	"os"
	"strings"
	"testing"

	"github.com/jogang0304/envcheck/internal"
)

func TestValidateTypes(t *testing.T) {
	testWithError := func(config *internal.Config, expectedErrorText string) {
		err := internal.ValidateTypes(config)
		if err == nil {
			t.Fatal("expected to get an error")
		}

		if !strings.Contains(err.Error(), expectedErrorText) {
			t.Fatalf("expected \"%s\" to contain \"%s\"", err.Error(), expectedErrorText)
		}
	}

	t.Run("Correct config and correct vars", func(t *testing.T) {
		config := internal.Config{
			Vars: []internal.VarEntry{
				{
					Name: "stringVar",
					Type: internal.StringType,
				},
				{
					Name: "intVar",
					Type: internal.IntType,
				},
				{
					Name: "boolVar",
					Type: internal.BoolType,
				},
				{
					Name: "floatVar",
					Type: internal.FloatType,
				},
				{
					Name: "floatVar2",
					Type: internal.FloatType,
				},
				{
					Name: "anyVar",
					Type: internal.AnyType,
				},
				{
					Name: "emptyVar",
					Type: internal.FloatType,
				},
			},
		}

		if err := os.Setenv("stringVar", "testString"); err != nil {
			t.Fatal("failed to set env stringVar")
		}

		if err := os.Setenv("intVar", "1234"); err != nil {
			t.Fatal("failed to set env intVar")
		}

		if err := os.Setenv("floatVar", "-27.917"); err != nil {
			t.Fatal("failed to set env floatVar")
		}

		if err := os.Setenv("floatVar2", "0"); err != nil {
			t.Fatal("failed to set env floatVar2")
		}

		if err := os.Setenv("boolVar", "false"); err != nil {
			t.Fatal("failed to set env boolVar")
		}

		if err := os.Setenv("anyVar", "{[***]}{.,.,~!@#$%^&*()}"); err != nil {
			t.Fatal("failled to set env anyVar")
		}

		err := internal.ValidateTypes(&config)
		if err != nil {
			t.Fatalf("unexpected error \"%v\"", err)
		}
	})

	t.Run("Incorrect var type", func(t *testing.T) {
		t.Run("Incorrect int", func(t *testing.T) {
			config := internal.Config{
				Vars: []internal.VarEntry{
					{
						Name: "intVar",
						Type: internal.IntType,
					},
				},
			}

			if err := os.Setenv("intVar", "123def45"); err != nil {
				t.Fatal("failed to set env intVar")
			}

			const expectedErrorText = "var intVar is not a valid int"

			testWithError(&config, expectedErrorText)
		})

		t.Run("Incorrect float", func(t *testing.T) {
			config := internal.Config{
				Vars: []internal.VarEntry{
					{
						Name: "floatVar",
						Type: internal.FloatType,
					},
				},
			}

			if err := os.Setenv("floatVar", "123def45"); err != nil {
				t.Fatal("failed to set env floatVar")
			}

			const expectedErrorText = "var floatVar is not a valid float"

			testWithError(&config, expectedErrorText)
		})

		t.Run("Incorrect bool", func(t *testing.T) {
			config := internal.Config{
				Vars: []internal.VarEntry{
					{
						Name: "boolVar",
						Type: internal.BoolType,
					},
				},
			}

			if err := os.Setenv("boolVar", "123def45"); err != nil {
				t.Fatal("failed to set env boolVar")
			}

			const expectedErrorText = "var boolVar is not a valid bool"

			testWithError(&config, expectedErrorText)
		})

		t.Run("Unsupported type", func(t *testing.T) {
			config := internal.Config{
				Vars: []internal.VarEntry{
					{
						Name: "strangeVar",
						Type: "strange",
					},
				},
			}

			if err := os.Setenv("strangeVar", "123def45"); err != nil {
				t.Fatal("failed to set env strangeVar")
			}

			const expectedErrorText = "var strangeVar has an unsupported type"

			testWithError(&config, expectedErrorText)
		})
	})
}
