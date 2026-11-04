package pick_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/xynova/roster-pick/internal/config"
	derrors "github.com/xynova/roster-pick/internal/errors"
	"github.com/xynova/roster-pick/internal/observability"
	"github.com/xynova/roster-pick/internal/pick"
)

func testScoresPath(t *testing.T) string {
	t.Helper()
	return filepath.Join("testdata", "scores.json")
}

func TestPickDeepinfraFirst(t *testing.T) {
	t.Parallel()
	p, err := pick.CreatePicker(config.Config{
		ScoresPath: testScoresPath(t),
		Limit:      5,
	}, observability.NopLogger{})
	if err != nil {
		t.Fatal(err)
	}
	res, err := p.Pick(context.Background(), pick.Request{
		Task: "Industrial cheap vision please",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Offers) == 0 {
		t.Fatal("no offers")
	}
	if res.Offers[0].HostID != "deepinfra" {
		t.Fatalf("first host=%q", res.Offers[0].HostID)
	}
}

func TestPickRetiredExcluded(t *testing.T) {
	t.Parallel()
	p, err := pick.CreatePicker(config.Config{
		ScoresPath: testScoresPath(t),
	}, observability.NopLogger{})
	if err != nil {
		t.Fatal(err)
	}
	res, err := p.Pick(context.Background(), pick.Request{Task: "cheap vision industrial"})
	if err != nil {
		t.Fatal(err)
	}
	for _, o := range res.Offers {
		if o.Name == "old-model" {
			t.Fatal("retired offer included")
		}
	}
}

func TestPickMissingScoresUnavailable(t *testing.T) {
	t.Parallel()
	_, err := pick.CreatePicker(config.Config{
		ScoresPath: filepath.Join(t.TempDir(), "missing.json"),
	}, observability.NopLogger{})
	if err == nil {
		t.Fatal("expected error")
	}
	if !derrors.IsCode(err, derrors.CodeUnavailable) {
		t.Fatalf("code=%v", err)
	}
}
