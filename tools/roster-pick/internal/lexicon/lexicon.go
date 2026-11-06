package lexicon

import (
	"regexp"
	"sort"
	"strings"

	derrors "github.com/xynova/roster-pick/internal/errors"
)

var tokenSplit = regexp.MustCompile(`[^a-z]+`)

var jobAliases = map[string]string{
	"chat":            "chat",
	"reason":          "reasoning",
	"reasoning":       "reasoning",
	"code":            "coding",
	"coding":          "coding",
	"vision":          "vision",
	"image":           "vision",
	"multimodal":      "vision",
	"cheap":           "cheap",
	"budget":          "cheap",
	"cost":            "cheap",
	"throughput":      "high-throughput",
	"context":         "long-context",
	"long-context":    "long-context",
	"high-throughput": "high-throughput",
	"classify":        "classify",
	"decision":        "classify",
	"decider":         "classify",
	"systemone":       "classify",
	"tts":             "tts",
	"speech":          "tts",
	"voice":           "tts",
}

var capacityAliases = map[string]string{
	"industrial": "industrial",
	"production": "industrial",
	"batch":      "industrial",
	"hobby":      "hobby",
	"dev":        "hobby",
	"local":      "local",
	"offline":    "local",
}

var kindAliases = map[string]string{
	"sovereign": "sovereign",
	"portable":  "sovereign",
	"provider":  "provider",
	"service":   "service",
	"cursor":    "service",
}

var allowedJobs = []string{
	"chat", "reasoning", "coding", "vision", "long-context", "cheap", "high-throughput", "classify", "tts",
}

// Match holds parsed task axes.
type Match struct {
	Jobs     []string
	Capacity string
	Kind     string
}

// ParseTask maps a natural-language task onto taxonomy jobs, capacity, and kind.
func ParseTask(task string) (Match, error) {
	tokens := tokenSplit.Split(strings.ToLower(task), -1)
	jobSet := map[string]struct{}{}
	var capacity, kind string
	for _, tok := range tokens {
		if tok == "" {
			continue
		}
		if j, ok := jobAliases[tok]; ok {
			jobSet[j] = struct{}{}
		}
		if c, ok := capacityAliases[tok]; ok {
			capacity = c
		}
		if k, ok := kindAliases[tok]; ok {
			kind = k
		}
	}
	if len(jobSet) == 0 {
		return Match{}, derrors.New(derrors.CodeInvalidArgument, "lexicon.ParseTask", "no job tokens; allowed: "+strings.Join(allowedJobs, ", "))
	}
	jobs := make([]string, 0, len(jobSet))
	for j := range jobSet {
		jobs = append(jobs, j)
	}
	sort.Strings(jobs)
	return Match{Jobs: jobs, Capacity: capacity, Kind: kind}, nil
}
