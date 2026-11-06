package card

import (
	"fmt"
	"path"
)

// AssignCategory returns the sticky sidebar band for fm (taxonomy.md).
func AssignCategory(fm map[string]any) string {
	jobs, _ := fm["jobs"].([]any)
	for _, j := range jobs {
		s, _ := j.(string)
		if s == "classify" {
			return "classify"
		}
	}
	for _, j := range jobs {
		s, _ := j.(string)
		if s == "tts" {
			return "tts"
		}
	}
	if str(fm["kind"]) == "service" {
		return "unpriced"
	}
	var cheapest *float64
	hosts, _ := fm["hosts"].([]any)
	for _, h := range hosts {
		hm, ok := h.(map[string]any)
		if !ok {
			continue
		}
		if _, isBool := hm["price_in"].(bool); isBool {
			continue
		}
		if _, isBool := hm["price_out"].(bool); isBool {
			continue
		}
		pin, okIn := asFloat(hm["price_in"])
		pout, okOut := asFloat(hm["price_out"])
		if !okIn || !okOut {
			continue
		}
		total := pin + pout
		if cheapest == nil || total < *cheapest {
			t := total
			cheapest = &t
		}
	}
	if cheapest == nil {
		return "unpriced"
	}
	if *cheapest < CheapMax {
		return "cheap"
	}
	if *cheapest < MidMax {
		return "mid"
	}
	return "spend"
}

// ValidateCard checks taxonomy allowlists on source frontmatter.
func ValidateCard(fm map[string]any, relPath string) error {
	var missing []string
	for _, k := range cardRequired {
		if _, ok := fm[k]; !ok {
			missing = append(missing, k)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("%s: missing required fields: %v", relPath, missing)
	}
	kind := str(fm["kind"])
	if _, ok := AllowedKinds[kind]; !ok {
		return fmt.Errorf("%s: invalid kind %q", relPath, kind)
	}
	cat := str(fm["category"])
	if _, ok := AllowedCategories[cat]; !ok {
		return fmt.Errorf("%s: invalid category %q", relPath, cat)
	}
	st := str(fm["status"])
	if _, ok := AllowedStatus[st]; !ok {
		return fmt.Errorf("%s: invalid status %q", relPath, st)
	}
	jobs, ok := fm["jobs"].([]any)
	if !ok || len(jobs) == 0 {
		return fmt.Errorf("%s: jobs must be a non-empty list", relPath)
	}
	for _, job := range jobs {
		j := str(job)
		if _, ok := AllowedJobs[j]; !ok {
			return fmt.Errorf("%s: invalid job %q", relPath, j)
		}
	}
	if kind == "service" && str(fm["service"]) == "" {
		return fmt.Errorf("%s: kind service requires service", relPath)
	}
	hosts, ok := fm["hosts"].([]any)
	if !ok || len(hosts) == 0 {
		return fmt.Errorf("%s: hosts must be a non-empty list", relPath)
	}
	seen := map[string]struct{}{}
	for i, host := range hosts {
		hm, ok := host.(map[string]any)
		if !ok {
			return fmt.Errorf("%s: hosts[%d] must be a map", relPath, i)
		}
		if _, hasKind := hm["kind"]; hasKind {
			return fmt.Errorf("%s: hosts[%d] must not contain kind", relPath, i)
		}
		hid := str(hm["id"])
		if _, ok := AllowedHostIDs[hid]; !ok {
			return fmt.Errorf("%s: hosts[%d] invalid id %q", relPath, i, hid)
		}
		if _, dup := seen[hid]; dup {
			return fmt.Errorf("%s: duplicate host id %q", relPath, hid)
		}
		seen[hid] = struct{}{}
		cap := str(hm["capacity"])
		if _, ok := AllowedCapacity[cap]; !ok {
			return fmt.Errorf("%s: hosts[%d] invalid capacity %q", relPath, i, cap)
		}
		if kind == "sovereign" || kind == "provider" {
			runtime := str(hm["runtime"])
			if _, ok := AllowedRuntimes[runtime]; !ok {
				return fmt.Errorf("%s: hosts[%d] requires runtime for %s cards", relPath, i, kind)
			}
		}
	}
	preferred := str(fm["preferred_host"])
	if preferred != "" {
		if _, ok := seen[preferred]; !ok {
			return fmt.Errorf("%s: preferred_host %q does not match a host id", relPath, preferred)
		}
	}
	_ = path.Clean(relPath)
	return nil
}

// ToJSON builds a catalog.json card object.
func ToJSON(fm map[string]any) map[string]any {
	out := map[string]any{
		"name":     fm["name"],
		"title":    fm["title"],
		"kind":     fm["kind"],
		"category": fm["category"],
		"jobs":     fm["jobs"],
		"status":   fm["status"],
		"lab":      fm["lab"],
		"hosts":    []any{},
	}
	for _, key := range []string{
		"preferred_host", "service", "aa_intelligence",
		"arena_elo", "aider_polyglot", "context_window", "links",
	} {
		if v, ok := fm[key]; ok && !isEmptyUnknown(v) {
			out[key] = unknownToJSON(v)
		}
	}
	hosts, _ := fm["hosts"].([]any)
	hout := make([]any, 0, len(hosts))
	for _, host := range hosts {
		hm, ok := host.(map[string]any)
		if !ok {
			continue
		}
		row := map[string]any{
			"id":       hm["id"],
			"capacity": hm["capacity"],
		}
		for _, key := range []string{
			"runtime", "runtime_model_id", "price_in", "price_out",
			"price_note", "limits_note", "links",
		} {
			if v, ok := hm[key]; ok && !isEmptyUnknown(v) {
				row[key] = unknownToJSON(v)
			}
		}
		hout = append(hout, row)
	}
	out["hosts"] = hout
	return out
}

func unknownToJSON(value any) any {
	if s, ok := value.(string); ok && s == "unknown" {
		return nil
	}
	switch t := value.(type) {
	case map[string]any:
		out := make(map[string]any, len(t))
		for k, v := range t {
			out[k] = unknownToJSON(v)
		}
		return out
	case []any:
		out := make([]any, len(t))
		for i, v := range t {
			out[i] = unknownToJSON(v)
		}
		return out
	default:
		return value
	}
}

func isEmptyUnknown(v any) bool {
	if v == nil {
		return true
	}
	s, ok := v.(string)
	return ok && (s == "" || s == "unknown")
}

func str(v any) string {
	s, _ := v.(string)
	return s
}

func asFloat(v any) (float64, bool) {
	switch t := v.(type) {
	case bool:
		return 0, false
	case int:
		return float64(t), true
	case int64:
		return float64(t), true
	case uint64:
		return float64(t), true
	case float32:
		return float64(t), true
	case float64:
		return t, true
	default:
		return 0, false
	}
}
