package internal

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
)

/*
These functions reads .env file specified by path f.
They iterate by lines and set os.env variables.
If the file is not found, it returns an error.
*/
func LoadDotenvFromFile(f string) error {
	var loadError error = nil

	fileBytes, err := os.ReadFile(f)
	if err != nil {
		return fmt.Errorf("failed to read .env\n\t %w", err)
	}

	// Read the file line by line
	scanner := bufio.NewScanner(bytes.NewReader(fileBytes))
	for scanner.Scan() {
		line := scanner.Text()
		// Skip empty lines and comments
		if line == "" || line[0] == '#' {
			continue
		}
		// Split the line into key and value
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			loadError = errors.Join(fmt.Errorf("invalid line in .env file: %s", line), loadError)
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// If value is covered in quotes, remove them
		if len(value) > 1 && value[0] == '"' && value[len(value)-1] == '"' {
			value = value[1 : len(value)-1]
		} else if len(value) > 1 && value[0] == '\'' && value[len(value)-1] == '\'' {
			value = value[1 : len(value)-1]
		}

		// Set the environment variable
		if err := os.Setenv(key, value); err != nil {
			loadError = errors.Join(
				fmt.Errorf("failed to set env var %s\n\t%w", key, err),
				loadError,
			)
		}
	}
	return loadError
}

/*
Detects in which directory go.mod is located and loads .env file from there.
*/
func LoadDotenv() error {
	var loadError error = nil

	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		loadError = errors.Join(
			errors.Join(errors.New("failed to get current working directory"), err),
			loadError,
		)
	}

	if err := LoadDotenvFromFile(cwd + "/.env"); err != nil {
		loadError = errors.Join(fmt.Errorf("failed to load .env file\n\t%w", err), loadError)
	}

	return loadError
}
