package repo

import (
	"os"
	"path/filepath"

	derrors "github.com/xynova/roster-mill/internal/errors"
)

// FindRoot walks up from start until it finds roster/taxonomy.md.
func FindRoot(start string) (string, error) {
	if start == "" {
		wd, err := os.Getwd()
		if err != nil {
			return "", derrors.Wrap(err, derrors.CodeFailed, "repo.FindRoot", "working directory")
		}
		start = wd
	}
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", derrors.Wrap(err, derrors.CodeFailed, "repo.FindRoot", "abs")
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "roster", "taxonomy.md")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", derrors.New(derrors.CodeNotFound, "repo.FindRoot", "roster/taxonomy.md not found from "+start)
		}
		dir = parent
	}
}
