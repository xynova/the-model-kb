package cli

import (
	"github.com/spf13/cobra"

	"github.com/xynova/roster-pick/internal/mcp"
)

func newMCPCmd() *cobra.Command {
	var scoresPath string
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "Start stdio MCP server (tool roster_pick)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return mcp.ServeStdio(mcp.ServerConfig{
				ScoresPath: scoresPath,
				Version:    Version,
			})
		},
	}
	cmd.Flags().StringVar(&scoresPath, "scores", "", "Path to scores.json (default: discovery)")
	return cmd
}
