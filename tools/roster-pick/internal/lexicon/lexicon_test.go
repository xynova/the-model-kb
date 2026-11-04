package lexicon_test

import (
	"testing"

	derrors "github.com/xynova/roster-pick/internal/errors"
	"github.com/xynova/roster-pick/internal/lexicon"
)

func TestLexiconIndustrialCheapVision(t *testing.T) {
	t.Parallel()
	m, err := lexicon.ParseTask("Industrial cheap vision please")
	if err != nil {
		t.Fatal(err)
	}
	if m.Capacity != "industrial" {
		t.Fatalf("capacity=%q", m.Capacity)
	}
	wantJobs := map[string]bool{"cheap": true, "vision": true}
	for _, j := range m.Jobs {
		if !wantJobs[j] {
			t.Fatalf("unexpected job %q in %v", j, m.Jobs)
		}
		delete(wantJobs, j)
	}
	if len(wantJobs) != 0 {
		t.Fatalf("missing jobs %v", wantJobs)
	}
}

func TestLexiconNoJob(t *testing.T) {
	t.Parallel()
	_, err := lexicon.ParseTask("hello")
	if err == nil {
		t.Fatal("expected error")
	}
	if !derrors.IsCode(err, derrors.CodeInvalidArgument) {
		t.Fatalf("code=%v", err)
	}
}
