package apply

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/xynova/roster-mill/internal/aa"
	"github.com/xynova/roster-mill/internal/card"
	derrors "github.com/xynova/roster-mill/internal/errors"
	"gopkg.in/yaml.v3"
)

var clearStatuses = map[string]struct{}{"no_score": {}, "unmatched": {}}

// Decision is a confirmed capacity change.
type Decision struct {
	Name     string
	HostID   string
	Capacity string
}

func sameScore(left, right any) bool {
	a := aa.YAMLScore(left)
	b := aa.YAMLScore(right)
	if a == nil || b == nil {
		return false
	}
	af, okA := asF(a)
	bf, okB := asF(b)
	return okA && okB && af == bf
}

func hasFiledScore(value any) bool {
	return aa.YAMLScore(value) != nil
}

func asF(v any) (float64, bool) {
	switch t := v.(type) {
	case int:
		return float64(t), true
	case float64:
		return t, true
	default:
		return 0, false
	}
}

func scoreNode(v any) *yaml.Node {
	s := aa.YAMLScore(v)
	switch t := s.(type) {
	case int:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: fmt.Sprintf("%d", t)}
	case float64:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!float", Value: fmt.Sprintf("%g", t)}
	default:
		return nil
	}
}

func floatNode(v float64) *yaml.Node {
	if v == float64(int(v)) {
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!float", Value: fmt.Sprintf("%.1f", v)}
	}
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!float", Value: fmt.Sprintf("%g", v)}
}

func isNumber(v any) bool {
	switch v.(type) {
	case int, int64, float32, float64, json.Number:
		return true
	default:
		return false
	}
}

// LoadDecisions reads capacity-decisions.yaml.
func LoadDecisions(path string) ([]Decision, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var data map[string]any
	if err := yaml.Unmarshal(raw, &data); err != nil {
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf("%s: expected mapping with decisions", path)
	}
	decisions, ok := data["decisions"]
	if !ok {
		return nil, fmt.Errorf("%s: missing decisions key", path)
	}
	list, ok := decisions.([]any)
	if !ok {
		return nil, fmt.Errorf("%s: decisions must be a list", path)
	}
	seen := map[string]struct{}{}
	var out []Decision
	for i, row := range list {
		m, ok := row.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("%s: decisions[%d] must be a map", path, i)
		}
		name := str(m["name"])
		hostID := str(m["host_id"])
		capacity := str(m["capacity"])
		if name == "" || hostID == "" {
			return nil, fmt.Errorf("%s: decisions[%d] needs name and host_id", path, i)
		}
		if _, ok := card.AllowedCapacity[capacity]; !ok {
			return nil, fmt.Errorf("%s: decisions[%d] invalid capacity %q", path, i, capacity)
		}
		key := name + "\x00" + hostID
		if _, dup := seen[key]; dup {
			return nil, fmt.Errorf("%s: duplicate decision (%s, %s)", path, name, hostID)
		}
		seen[key] = struct{}{}
		out = append(out, Decision{Name: name, HostID: hostID, Capacity: capacity})
	}
	return out, nil
}

func decisionMap(decisions []Decision) map[string]string {
	out := map[string]string{}
	for _, d := range decisions {
		out[d.Name+"\x00"+d.HostID] = d.Capacity
	}
	return out
}

func rowList(review map[string]any) []map[string]any {
	raw, _ := review["rows"].([]any)
	if raw == nil {
		if typed, ok := review["rows"].([]map[string]any); ok {
			return typed
		}
		return nil
	}
	out := make([]map[string]any, 0, len(raw))
	for _, item := range raw {
		if m, ok := item.(map[string]any); ok {
			out = append(out, m)
		}
	}
	return out
}

// ApplyMine patches host notes, prices, and confirmed capacity onto cards.
func ApplyMine(review map[string]any, decisions []Decision, cardsRoot string, now time.Time) (int, []string) {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	today := now.UTC().Format("2006-01-02")
	dmap := decisionMap(decisions)
	var errors []string
	for _, row := range rowList(review) {
		if str(row["status"]) != "ok" {
			continue
		}
		name := str(row["name"])
		hostID := str(row["host_id"])
		if name == "" || hostID == "" {
			errors = append(errors, fmt.Sprintf("review row missing name or host_id: %v", row))
			continue
		}
		path, f, err := card.FindByName(cardsRoot, name)
		if err != nil {
			errors = append(errors, err.Error())
			continue
		}
		if path == "" || f == nil {
			errors = append(errors, fmt.Sprintf("no card for name %q", name))
			continue
		}
		fm := f.AsMap()
		if fm == nil {
			errors = append(errors, path+": invalid frontmatter")
			continue
		}
		hostsNode := card.Lookup(card.Mapping(f.Front), "hosts")
		_, hostNode := card.HostMapping(fm, hostsNode, hostID)
		if hostNode == nil || hostNode.Kind != yaml.MappingNode {
			errors = append(errors, fmt.Sprintf("%s: no host id %q", name, hostID))
			continue
		}
		if evidence := str(row["evidence"]); evidence != "" {
			card.Set(hostNode, "limits_note", card.ScalarString(evidence))
		}
		if isNumber(row["price_in"]) && isNumber(row["price_out"]) {
			pin, _ := asF(aa.YAMLScore(row["price_in"]))
			pout, _ := asF(aa.YAMLScore(row["price_out"]))
			card.Set(hostNode, "price_in", floatNode(pin))
			card.Set(hostNode, "price_out", floatNode(pout))
			if evidence := str(row["evidence"]); evidence != "" {
				card.Set(hostNode, "price_note", card.ScalarString(evidence))
			}
		}
		if cap, ok := dmap[name+"\x00"+hostID]; ok {
			card.Set(hostNode, "capacity", card.ScalarString(cap))
			row["capacity_applied"] = true
		}
		card.Set(card.Mapping(f.Front), "refreshed", card.ScalarQuoted(today))
		if err := card.WriteCard(path, f.Front, f.Body); err != nil {
			errors = append(errors, err.Error())
		}
	}
	code := 0
	if len(errors) > 0 {
		code = 1
	}
	return code, errors
}

// ApplyAA files or clears aa_intelligence on cards from the AA review.
func ApplyAA(review map[string]any, cardsRoot string, now time.Time) (int, []string) {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	if str(review["status"]) != "ok" {
		return 0, nil
	}
	today := now.UTC().Format("2006-01-02")
	var errors []string
	for _, row := range rowList(review) {
		name := str(row["name"])
		if name == "" {
			continue
		}
		status := str(row["status"])
		if status == "ok" {
			score := aa.YAMLScore(row["aa_intelligence"])
			slug := str(row["aa_slug"])
			if score == nil {
				errors = append(errors, fmt.Sprintf("review row missing name or score: %v", row))
				continue
			}
			path, f, err := card.FindByName(cardsRoot, name)
			if err != nil {
				errors = append(errors, err.Error())
				continue
			}
			if path == "" || f == nil {
				errors = append(errors, fmt.Sprintf("no card for name %q", name))
				continue
			}
			fm := f.AsMap()
			if fm == nil {
				errors = append(errors, path+": invalid frontmatter")
				continue
			}
			previous := fm["aa_intelligence"]
			row["previous"] = previous
			if sameScore(previous, score) {
				row["applied"] = false
				continue
			}
			m := card.Mapping(f.Front)
			card.Set(m, "aa_intelligence", scoreNode(score))
			links := card.Lookup(m, "links")
			if links == nil || links.Kind != yaml.MappingNode {
				links = &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
				card.Set(m, "links", links)
			}
			if slug != "" {
				card.Set(links, "artificial_analysis", card.ScalarString(aa.ModelPageURL(slug)))
			}
			card.Set(m, "refreshed", card.ScalarQuoted(today))
			if err := card.WriteCard(path, f.Front, f.Body); err != nil {
				errors = append(errors, err.Error())
				continue
			}
			row["applied"] = true
			continue
		}
		if _, ok := clearStatuses[status]; !ok {
			continue
		}
		path, f, err := card.FindByName(cardsRoot, name)
		if err != nil {
			errors = append(errors, err.Error())
			continue
		}
		if path == "" || f == nil {
			errors = append(errors, fmt.Sprintf("no card for name %q", name))
			continue
		}
		fm := f.AsMap()
		if fm == nil {
			errors = append(errors, path+": invalid frontmatter")
			continue
		}
		previous := fm["aa_intelligence"]
		row["previous"] = previous
		if !hasFiledScore(previous) {
			row["applied"] = false
			continue
		}
		m := card.Mapping(f.Front)
		card.Set(m, "aa_intelligence", card.ScalarString("unknown"))
		card.Set(m, "refreshed", card.ScalarQuoted(today))
		if err := card.WriteCard(path, f.Front, f.Body); err != nil {
			errors = append(errors, err.Error())
			continue
		}
		row["applied"] = true
	}
	code := 0
	if len(errors) > 0 {
		code = 1
	}
	return code, errors
}

func printErrors(errs []string) {
	for _, msg := range errs {
		fmt.Fprintln(os.Stderr, msg)
	}
}

func loadReview(path string) (map[string]any, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var review map[string]any
	if err := json.Unmarshal(raw, &review); err != nil {
		return nil, err
	}
	return review, nil
}

func rewriteReview(path string, review map[string]any) error {
	return aa.WriteJSON(path, review)
}

// RunMine is the CLI apply-mine command.
func RunMine(repoRoot string) error {
	reviewPath := filepath.Join(repoRoot, "roster", "retrieval", "limits-review.json")
	decisionsPath := filepath.Join(repoRoot, "roster", "retrieval", "capacity-decisions.yaml")
	cardsRoot := filepath.Join(repoRoot, "roster")
	if _, err := os.Stat(reviewPath); err != nil {
		fmt.Fprintf(os.Stderr, "review missing: %s\n", reviewPath)
		return derrors.Wrap(err, derrors.CodeNotFound, "apply-mine", "review missing")
	}
	decisions, err := LoadDecisions(decisionsPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return derrors.Wrap(err, derrors.CodeInvalidArgument, "apply-mine", "capacity decisions")
	}
	review, err := loadReview(reviewPath)
	if err != nil {
		return derrors.Wrap(err, derrors.CodeFailed, "apply-mine", "parse review")
	}
	code, errs := ApplyMine(review, decisions, cardsRoot, time.Now().UTC())
	if werr := rewriteReview(reviewPath, review); werr != nil {
		return werr
	}
	printErrors(errs)
	if code != 0 {
		return derrors.New(derrors.CodeFailed, "apply-mine", "card apply errors")
	}
	return nil
}

// RunAA is the CLI apply-aa command.
func RunAA(repoRoot string) error {
	reviewPath := filepath.Join(repoRoot, "roster", "retrieval", "aa-intelligence-review.json")
	cardsRoot := filepath.Join(repoRoot, "roster")
	if _, err := os.Stat(reviewPath); err != nil {
		fmt.Fprintf(os.Stderr, "AA review missing: %s\n", reviewPath)
		return nil
	}
	review, err := loadReview(reviewPath)
	if err != nil {
		return derrors.Wrap(err, derrors.CodeFailed, "apply-aa", "parse review")
	}
	code, errs := ApplyAA(review, cardsRoot, time.Now().UTC())
	if werr := rewriteReview(reviewPath, review); werr != nil {
		return werr
	}
	printErrors(errs)
	if code != 0 {
		return derrors.New(derrors.CodeFailed, "apply-aa", "card apply errors")
	}
	return nil
}

func str(v any) string {
	s, _ := v.(string)
	return s
}
