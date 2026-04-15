package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/Gauravism2017/cctl/internal/profile"
)

var removeFlag bool

var initCmd = &cobra.Command{
	Use:   "init <profile-name>",
	Short: "Bind a profile to the current directory",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profileName := args[0]
		claudeDir := filepath.Join(".claude")
		profileFile := filepath.Join(claudeDir, "profile")

		if err := os.MkdirAll(claudeDir, 0o755); err != nil {
			return fmt.Errorf("creating .claude directory: %w", err)
		}

		existing := readProfileLines(profileFile)

		if removeFlag {
			updated := removeFromLines(existing, profileName)
			if len(updated) == len(existing) {
				fmt.Printf("Profile %q is not bound to this directory\n", profileName)
				return nil
			}
			if err := writeProfileLines(profileFile, updated); err != nil {
				return err
			}
			fmt.Printf("Unbound profile %q from this directory\n", profileName)
			return nil
		}

		if contains(existing, profileName) {
			fmt.Printf("Profile %q is already bound to this directory\n", profileName)
			return nil
		}

		if _, err := profile.Get(profileName, paths); err != nil {
			if err := profile.Save(profileName, paths); err != nil {
				return fmt.Errorf("creating profile: %w", err)
			}
			fmt.Printf("Profile %q created from current state\n", profileName)
		}

		updated := append(existing, profileName)
		if err := writeProfileLines(profileFile, updated); err != nil {
			return err
		}

		fmt.Printf("Bound profile %q to this directory\n", profileName)
		return nil
	},
}

func readProfileLines(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var lines []string
	for _, line := range strings.Split(string(data), "\n") {
		name := strings.TrimSpace(line)
		if name != "" {
			lines = append(lines, name)
		}
	}
	return lines
}

func writeProfileLines(path string, names []string) error {
	content := strings.Join(names, "\n") + "\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("writing profile file: %w", err)
	}
	return nil
}

func removeFromLines(lines []string, name string) []string {
	var result []string
	for _, l := range lines {
		if l != name {
			result = append(result, l)
		}
	}
	return result
}

func contains(lines []string, name string) bool {
	for _, l := range lines {
		if l == name {
			return true
		}
	}
	return false
}

func init() {
	initCmd.Flags().BoolVarP(&removeFlag, "remove", "r", false, "unbind profile from this directory")
}
