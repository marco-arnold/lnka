package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Config holds the application configuration
type Config struct {
	SourceDir string
	TargetDir string
	Title     string
}

// Load loads configuration from cobra command
func Load(cmd *cobra.Command, args []string) (*Config, error) {
	cfg := &Config{}

	// Get positional arguments (source and target)
	if len(args) >= 1 {
		cfg.SourceDir = args[0]
	}

	if len(args) >= 2 {
		cfg.TargetDir = args[1]
	}

	// Get flags
	var err error
	cfg.Title, err = cmd.Flags().GetString("title")
	if err != nil {
		return nil, fmt.Errorf("failed to get title flag: %w", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// Check if both directories are provided
	if c.SourceDir == "" {
		return errors.New("source directory not specified: provide as first argument")
	}

	if c.TargetDir == "" {
		return errors.New("target directory not specified: provide as second argument")
	}

	// Check if directories exist
	if err := checkDirExists(c.SourceDir); err != nil {
		return fmt.Errorf("source directory: %w", err)
	}

	if err := checkDirExists(c.TargetDir); err != nil {
		return fmt.Errorf("target directory: %w", err)
	}

	return nil
}

// checkDirExists verifies that a directory exists and is accessible
func checkDirExists(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s does not exist", path)
		}
		return fmt.Errorf("cannot access %s: %w", path, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}

	return nil
}
