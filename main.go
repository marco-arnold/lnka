package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/marco-arnold/lnka/internal/config"
	"github.com/marco-arnold/lnka/internal/filesystem"
	"github.com/marco-arnold/lnka/internal/ui"
	"github.com/spf13/cobra"
)

// Version information (set by goreleaser via ldflags)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "lnka SOURCE TARGET",
	Short: "Manage symlinks between source and target directories",
	Long: `lnka is a CLI tool for managing symlinks between a source directory
and a target directory using an interactive Terminal UI.`,
	Args: func(cmd *cobra.Command, args []string) error {
		// Allow version flag without args
		if versionFlag, _ := cmd.Flags().GetBool("version"); versionFlag {
			return nil
		}
		return cobra.ExactArgs(2)(cmd, args)
	},
	RunE: run,
}

func init() {
	// Define flags with environment variable fallback and shorthands
	titleDefault := os.Getenv("LNKA_TITLE")
	rootCmd.Flags().StringP("title", "t", titleDefault, "Title to display in UI (env: LNKA_TITLE)")

	// Add version flag
	rootCmd.Flags().BoolP("version", "v", false, "Print version information")

	// Add debug flag
	rootCmd.Flags().StringP("debug", "d", "", "Enable debug logging to specified file (e.g., debug.log)")
}

func printVersion() {
	fmt.Printf("lnka %s\n", version)
	fmt.Printf("  commit: %s\n", commit)
	fmt.Printf("  built at: %s\n", date)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		// Cobra already prints the error, just exit
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	// Check for version flag
	if versionFlag, _ := cmd.Flags().GetBool("version"); versionFlag {
		printVersion()
		return nil
	}

	// Setup debug logging if debug flag is set
	debugFile, _ := cmd.Flags().GetString("debug")
	if debugFile != "" {
		// Remove existing debug file to start fresh
		_ = os.Remove(debugFile)

		f, err := tea.LogToFile(debugFile, "lnka")
		if err != nil {
			return fmt.Errorf("failed to setup debug logging: %w", err)
		}
		defer f.Close()
	}

	// Load configuration
	cfg, err := config.Load(cmd, args)
	if err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	// Check for orphaned symlinks
	orphaned, err := filesystem.ValidateSymlinks(cfg.SourceDir, cfg.TargetDir)
	if err != nil {
		return fmt.Errorf("failed to validate symlinks: %w", err)
	}

	// If there are orphaned symlinks, ask user if they want to clean them
	if len(orphaned) > 0 {
		fmt.Printf("Found %d orphaned symlink(s):\n", len(orphaned))
		for _, name := range orphaned {
			fmt.Printf("  - %s\n", name)
		}
		fmt.Println()

		confirmed, err := ui.ShowConfirmation("Do you want to clean these orphaned symlinks?")
		if err != nil {
			if strings.Contains(err.Error(), "user aborted") {
				os.Exit(1)
			}
			return err
		}

		if confirmed {
			if err := filesystem.CleanOrphanedSymlinks(cfg.TargetDir, orphaned); err != nil {
				return fmt.Errorf("failed to clean orphaned symlinks: %w", err)
			}
			fmt.Printf("Cleaned %d orphaned symlink(s)\n\n", len(orphaned))
		}
	}

	// Show multi-select UI (loads files asynchronously in Init())
	selectedFiles, err := ui.ShowFileSelect(cfg.SourceDir, cfg.TargetDir, cfg.Title)
	if err != nil {
		if strings.Contains(err.Error(), "user aborted") {
			os.Exit(1)
		}
		return err
	}

	// Apply changes
	if err := filesystem.ApplyChanges(cfg.SourceDir, cfg.TargetDir, selectedFiles); err != nil {
		return fmt.Errorf("failed to apply changes: %w", err)
	}

	return nil
}
