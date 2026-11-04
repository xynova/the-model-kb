package config

import (
	"os"
	"path/filepath"

	derrors "github.com/xynova/roster-pick/internal/errors"
)

// Config holds paths and defaults for the picker.
type Config struct {
	ScoresPath string
	Limit      int
}

// CreateConfig resolves scores path from ROSTER_SCORES or walk-up discovery.
func CreateConfig() (Config, error) {
	path := os.Getenv("ROSTER_SCORES")
	if path == "" {
		var err error
		path, err = discoverScoresPath()
		if err != nil {
			return Config{}, err
		}
		if path == "" {
			return Config{}, derrors.New(derrors.CodeUnavailable, "config.CreateConfig", "scores file missing; run make -C dbt build")
		}
	}
	return Config{
		ScoresPath: path,
		Limit:      5,
	}, nil
}

func discoverScoresPath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for i := 0; i < 8; i++ {
		candidate := filepath.Join(dir, "roster", "scores.json")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", nil
}
