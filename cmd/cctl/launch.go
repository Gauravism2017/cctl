package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/Gauravism2017/cctl/internal/profile"
)

var profileFlag string

var launchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Load bound profile and launch claude",
	RunE: func(cmd *cobra.Command, args []string) error {
		profileName, err := resolveProfile()
		if err != nil {
			return err
		}

		applied, skipped, errs := profile.Load(profileName, paths)
		if len(errs) > 0 {
			for _, e := range errs {
				fmt.Fprintf(os.Stderr, "Error: %v\n", e)
			}
			return fmt.Errorf("failed to load profile %q", profileName)
		}

		fmt.Printf("Profile %q loaded: %d changes", profileName, applied)
		if skipped > 0 {
			fmt.Printf(", %d skipped", skipped)
		}
		fmt.Println()

		claudeBin, err := exec.LookPath("claude")
		if err != nil {
			return fmt.Errorf("claude not found in PATH")
		}

		return syscall.Exec(claudeBin, []string{"claude"}, os.Environ())
	},
}

func resolveProfile() (string, error) {
	if profileFlag != "" {
		return profileFlag, nil
	}

	names := readProfileLines(".claude/profile")
	if len(names) == 0 {
		return "", fmt.Errorf("no .claude/profile found — use 'cctl launch -p <name>' or 'cctl init <profile>' first")
	}

	// filter to profiles that actually exist
	var valid []string
	for _, name := range names {
		if _, err := profile.Get(name, paths); err == nil {
			valid = append(valid, name)
		}
	}

	switch len(valid) {
	case 0:
		return "", fmt.Errorf("no valid profiles found in .claude/profile (checked: %s)", strings.Join(names, ", "))
	case 1:
		return valid[0], nil
	default:
		return runPicker(valid)
	}
}

func init() {
	launchCmd.Flags().StringVarP(&profileFlag, "profile", "p", "", "profile to load (overrides .claude/profile)")
}
