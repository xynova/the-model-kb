package errors_test

import (
	"errors"
	"strings"
	"testing"

	derrors "github.com/xynova/library-intake/internal/errors"
)

func TestErrorIncludesSortedFields(t *testing.T) {
	t.Parallel()
	err := derrors.Wrap(errors.New("exit status 1"), derrors.CodeFailed, "execx.Run", "command failed").
		With("binary", "yt-dlp").
		With("stderr", "HTTP Error 403: Forbidden")
	got := err.Error()
	if !strings.Contains(got, "binary=yt-dlp") {
		t.Fatalf("missing binary field: %q", got)
	}
	if !strings.Contains(got, "stderr=HTTP Error 403: Forbidden") {
		t.Fatalf("missing stderr field: %q", got)
	}
	if !strings.Contains(got, "command failed: exit status 1") {
		t.Fatalf("missing cause: %q", got)
	}
	binIdx := strings.Index(got, "binary=")
	stderrIdx := strings.Index(got, "stderr=")
	if binIdx < 0 || stderrIdx < 0 || binIdx > stderrIdx {
		t.Fatalf("fields not sorted: %q", got)
	}
}
