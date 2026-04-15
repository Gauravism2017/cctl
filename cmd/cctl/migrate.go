package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/Gauravism2017/cctl/internal/migrate"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate from vault to symlink store",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !paths.NeedsMigration() {
			fmt.Println("Already using symlink store")
			return nil
		}
		result, err := migrate.Run(paths)
		if err != nil {
			return fmt.Errorf("migration error: %w", err)
		}
		fmt.Printf("Migrated to symlink store: %d skills, %d rules, %d agents\n",
			result.Skills, result.Rules, result.Agents)
		return nil
	},
}
