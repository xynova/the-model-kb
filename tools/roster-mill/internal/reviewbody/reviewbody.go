package reviewbody

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Markdown builds the roster-mine PR body from review JSON files.
func Markdown(limitsPath, aaPath string) (string, error) {
	raw, err := os.ReadFile(limitsPath)
	if err != nil {
		return "", err
	}
	var data map[string]any
	if err := json.Unmarshal(raw, &data); err != nil {
		return "", err
	}
	var b strings.Builder
	b.WriteString("## Roster mine review\n\n")
	rows, _ := data["rows"].([]any)
	for _, item := range rows {
		row, _ := item.(map[string]any)
		fmt.Fprintf(&b, "- **%v** / `%v`\n", row["name"], row["host_id"])
		u := row["url"]
		if u == nil || fmt.Sprint(u) == "" {
			u = "(skipped)"
		}
		fmt.Fprintf(&b, "  - url: %v\n", u)
		ev := row["evidence"]
		if ev == nil || fmt.Sprint(ev) == "" {
			ev = "(none)"
		}
		fmt.Fprintf(&b, "  - evidence: %v\n", ev)
		fmt.Fprintf(&b, "  - suggested_capacity: %v\n", row["suggested_capacity"])
		fmt.Fprintf(&b, "  - current_capacity: %v\n", row["current_capacity"])
		fmt.Fprintf(&b, "  - capacity_applied: %v\n\n", row["capacity_applied"])
	}
	if aaRaw, err := os.ReadFile(aaPath); err == nil {
		var aa map[string]any
		if json.Unmarshal(aaRaw, &aa) == nil {
			b.WriteString("## AA Intelligence Index\n\n")
			fmt.Fprintf(&b, "status: %v\n", aa["status"])
			fmt.Fprintf(&b, "index version: %v\n\n", aa["intelligence_index_version"])
			aaRows, _ := aa["rows"].([]any)
			for _, item := range aaRows {
				row, _ := item.(map[string]any)
				effort := row["effort"]
				if effort == nil || fmt.Sprint(effort) == "" {
					effort = "-"
				}
				fmt.Fprintf(&b, "- **%v**: %v score=%v effort=%v applied=%v\n",
					row["name"], row["status"], row["aa_intelligence"], effort, row["applied"])
			}
			b.WriteString("\n")
		}
	}
	return b.String(), nil
}
