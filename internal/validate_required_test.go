package internal_test

import (
	"os"
	"strings"
	"testing"

	"github.com/jogang0304/envcheck/internal"
)

func TestValidateRequired(t *testing.T) {
	t.Run("All vars have Required field", func(t *testing.T) {
		config := internal.Config{
			Vars: []internal.VarEntry{
				{
					Name:     "firstVar",
					Required: false,
				},
				{
					Name:     "secondVar",
					Required: true,
				},
			},
		}

		t.Run("Required vars are set", func(t *testing.T) {
			if err := os.Setenv("secondVar", "123"); err != nil {
				t.Fatal("failed to set env secondVar")
			}

			err := internal.ValidateRequired(&config)
			if err != nil {
				t.Fatalf("unexpected error \"%v\"", err)
			}
		})

		t.Run("Required vars are not set", func(t *testing.T) {
			if err := os.Unsetenv("secondVar"); err != nil {
				t.Fatal("failed to unset env secondVar")
			}

			const expectedErrorText = "required var secondVar is not set"

			err := internal.ValidateRequired(&config)
			if err == nil {
				t.Fatal("expected to get an error")
			}

			if !strings.Contains(err.Error(), expectedErrorText) {
				t.Fatalf("expected \"%s\" to contain %s", err.Error(), expectedErrorText)
			}
		})
	})

	t.Run("Some vars do not have Required field", func(t *testing.T) {
		config := internal.Config{
			Vars: []internal.VarEntry{
				{
					Name: "firstVar",
				},
				{
					Name:     "secondVar",
					Required: true,
				},
			},
		}

		t.Run("Required vars are set", func(t *testing.T) {
			if err := os.Setenv("secondVar", "123"); err != nil {
				t.Fatal("failed to set env secondVar")
			}

			if err := os.Unsetenv("firstVar"); err != nil {
				t.Fatal("failed to unset env firstVar")
			}

			err := internal.ValidateRequired(&config)
			if err != nil {
				t.Fatalf("unexpected error \"%v\"", err)
			}
		})
	})
}
