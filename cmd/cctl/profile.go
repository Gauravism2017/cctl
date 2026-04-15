package main

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/Gauravism2017/cctl/internal/profile"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage profiles",
}

var profileSaveCmd = &cobra.Command{
	Use:   "save <name>",
	Short: "Save current state as a profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := profile.Save(args[0], paths); err != nil {
			return err
		}
		fmt.Printf("Profile %q saved\n", args[0])
		return nil
	},
}

var profileLoadCmd = &cobra.Command{
	Use:   "load <name>",
	Short: "Load a profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		applied, skipped, errs := profile.Load(args[0], paths)
		if len(errs) > 0 {
			return fmt.Errorf("loading profile %q: %w", args[0], errors.Join(errs...))
		}
		fmt.Printf("Profile %q loaded: %d changes applied", args[0], applied)
		if skipped > 0 {
			fmt.Printf(", %d items skipped (no longer on disk)", skipped)
		}
		fmt.Println()
		return nil
	},
}

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all profiles",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := profile.List(paths)
		if err != nil {
			return err
		}
		if len(entries) == 0 {
			fmt.Println("No profiles saved yet")
			return nil
		}
		for _, e := range entries {
			fmt.Printf("  %-20s  %s\n", e.Name, e.Created)
		}
		return nil
	},
}

var profileDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := profile.Delete(args[0], paths); err != nil {
			return err
		}
		fmt.Printf("Profile %q deleted\n", args[0])
		return nil
	},
}

func init() {
	profileCmd.AddCommand(profileSaveCmd, profileLoadCmd, profileListCmd, profileDeleteCmd)
}
