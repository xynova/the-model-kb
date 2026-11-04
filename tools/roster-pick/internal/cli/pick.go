package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/xynova/roster-pick/internal/config"
	derrors "github.com/xynova/roster-pick/internal/errors"
	"github.com/xynova/roster-pick/internal/observability"
	"github.com/xynova/roster-pick/internal/pick"
)

func newPickCmd() *cobra.Command {
	var scoresPath string
	var asJSON bool
	var task string
	cmd := &cobra.Command{
		Use:   "pick",
		Short: "Rank offers for a natural-language task",
		RunE: func(cmd *cobra.Command, args []string) error {
			if task == "" {
				return derrors.New(derrors.CodeInvalidArgument, "cli.pick", "--task is required")
			}
			cfg := config.Config{Limit: 5}
			if scoresPath != "" {
				cfg.ScoresPath = scoresPath
			} else {
				c, err := config.CreateConfig()
				if err != nil {
					return writePickError(asJSON, err)
				}
				cfg = c
			}
			p, err := pick.CreatePicker(cfg, observability.NopLogger{})
			if err != nil {
				return writePickError(asJSON, err)
			}
			res, err := p.Pick(cmd.Context(), pick.Request{Task: task})
			if err != nil {
				if derrors.IsCode(err, derrors.CodeNotFound) {
					if wErr := writePickResult(asJSON, res); wErr != nil {
						return wErr
					}
					return err
				}
				return writePickError(asJSON, err)
			}
			return writePickResult(asJSON, res)
		},
	}
	cmd.Flags().StringVar(&task, "task", "", "Natural-language task description")
	cmd.Flags().StringVar(&scoresPath, "scores", "", "Path to scores.json (default: discovery)")
	cmd.Flags().BoolVar(&asJSON, "json", false, "Emit JSON")
	return cmd
}

func writePickResult(asJSON bool, res pick.Result) error {
	if asJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(res)
	}
	writeln(os.Stdout, "Matched jobs")
	writeln(os.Stdout, fmt.Sprintf("  %v", res.MatchedJobs))
	if res.MatchedCapacity != "" {
		writeln(os.Stdout, "Matched capacity")
		writeln(os.Stdout, fmt.Sprintf("  %s", res.MatchedCapacity))
	}
	writeln(os.Stdout, "Offers")
	for i, o := range res.Offers {
		cost := "null"
		if o.CostPerM != nil {
			cost = fmt.Sprintf("%g", *o.CostPerM)
		}
		writeln(os.Stdout, fmt.Sprintf("  %d. %s @ %s cost_per_m=%s", i+1, o.Name, o.HostID, cost))
	}
	return nil
}

func writePickError(asJSON bool, err error) error {
	if asJSON {
		code := derrors.CodeFailed
		var de *derrors.Error
		if errors.As(err, &de) {
			code = de.Code
		}
		payload := map[string]any{
			"error": map[string]string{
				"code":    string(code),
				"message": err.Error(),
			},
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(payload)
	}
	return err
}
