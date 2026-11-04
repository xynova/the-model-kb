package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// Version is set at link time; default for go run.
var Version = "dev"

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           "roster-pick",
		Short:         "Rank Model Picks host offers for a task",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			printAgentGuide(os.Stdout)
			return nil
		},
	}
	root.AddCommand(newVersionCmd(), newInspectCmd(), newPickCmd(), newMCPCmd())
	return root
}

func printAgentGuide(w io.Writer) {
	writeln(w, fmt.Sprintf("roster-pick %s", Version))
	writeln(w, "Rank Model Picks host offers for fit, cost, and host limits.")
	writeln(w, "")
	writeln(w, "Documentation for agents")
	writeln(w, "  tools/roster-pick/AGENTS.md")
	writeln(w, "  .cursor/skills/roster-pick/SKILL.md")
	writeln(w, "")
	writeln(w, "Read-only")
	writeln(w, "  inspect, pick, version, help")
	writeln(w, "")
	writeln(w, "Execute")
	writeln(w, "  mcp   Start stdio MCP (tool roster_pick)")
	writeln(w, "")
	writeln(w, "Full command list: roster-pick help")
}

func writeln(w io.Writer, line string) {
	_, _ = fmt.Fprintln(w, line)
}
