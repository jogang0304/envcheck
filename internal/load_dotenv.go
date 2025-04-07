package internal

import (
	"bufio"
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
	// Open the file
	file, err := os.Open(f)
	if err != nil {
		return err
	}
	defer file.Close()
	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Skip empty lines and comments
		if line == "" || line[0] == '#' {
			continue
		}
		// Split the line into key and value
		parts := strings.SplitN(line, "=", -1)
		if len(parts) != 2 {
			return fmt.Errorf("invalid line in .env file: %s", line)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		//If value is covered in quotes, remove them
		if len(value) > 1 && value[0] == '"' && value[len(value)-1] == '"' {
			value = value[1 : len(value)-1]
		} else if len(value) > 1 && value[0] == '\'' && value[len(value)-1] == '\'' {
			value = value[1 : len(value)-1]
		}

		// Set the environment variable
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("failed to set env var %s: %w", key, err)
		}
	}
	return nil
}

/*
Detects in which directory go.mod is located and loads .env file from there.
*/
func LoadDotenv() error {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	if err := LoadDotenvFromFile(cwd + "/.env"); err != nil {
		return fmt.Errorf("failed to load .env file: %w", err)
	}

	return nil
}
