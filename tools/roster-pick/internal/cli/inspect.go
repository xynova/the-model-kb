package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/xynova/roster-pick/internal/config"
	"github.com/xynova/roster-pick/internal/scores"
)

func newInspectCmd() *cobra.Command {
	var scoresPath string
	var asJSON bool
	cmd := &cobra.Command{
		Use:   "inspect",
		Short: "Show loaded scores metadata",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := scoresPath
			if path == "" {
				cfg, err := config.CreateConfig()
				if err != nil {
					return err
				}
				path = cfg.ScoresPath
			}
			doc, err := scores.Load(path)
			if err != nil {
				return err
			}
			if asJSON {
				payload := map[string]any{
					"scores_path":   path,
					"scores_schema": doc.ScoresSchema,
					"generated_at":  doc.GeneratedAt,
					"offer_count":   len(doc.Offers),
				}
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(payload)
			}
			writeln(os.Stdout, "Scores")
			writeln(os.Stdout, fmt.Sprintf("  Path:          %s", path))
			writeln(os.Stdout, fmt.Sprintf("  Schema:        %d", doc.ScoresSchema))
			writeln(os.Stdout, fmt.Sprintf("  Generated:     %s", doc.GeneratedAt))
			writeln(os.Stdout, fmt.Sprintf("  Offer count:   %d", len(doc.Offers)))
			return nil
		},
	}
	cmd.Flags().StringVar(&scoresPath, "scores", "", "Path to scores.json (default: discovery)")
	cmd.Flags().BoolVar(&asJSON, "json", false, "Emit JSON")
	return cmd
}
