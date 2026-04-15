package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/Gauravism2017/cctl/internal/config"
	"github.com/Gauravism2017/cctl/internal/migrate"
	"github.com/Gauravism2017/cctl/internal/scanner"
	"github.com/Gauravism2017/cctl/internal/tui"
)

var paths config.Paths

var rootCmd = &cobra.Command{
	Use:           "cctl",
	Short:         "Manage Claude Code context — skills, rules, agents, plugins, MCP servers, and profiles",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		paths, err = config.NewPaths()
		if err != nil {
			return err
		}

		if paths.NeedsMigration() {
			result, err := migrate.Run(paths)
			if err != nil {
				return fmt.Errorf("migration error: %w", err)
			}
			total := result.Skills + result.Rules + result.Agents
			if total > 0 {
				fmt.Printf("Migrated to symlink store: %d skills, %d rules, %d agents\n",
					result.Skills, result.Rules, result.Agents)
			}
		}

		return paths.EnsureDirs()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		skills := scanner.ScanSkills(paths)
		rules := scanner.ScanRules(paths)
		agents := scanner.ScanAgents(paths)
		plugins := scanner.ScanPlugins(paths)
		mcpServers := scanner.ScanMCPServers(paths)

		app := tui.NewApp(paths, skills, rules, agents, plugins, mcpServers)

		p := tea.NewProgram(app, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("TUI error: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd, initCmd, launchCmd, profileCmd)
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
}
