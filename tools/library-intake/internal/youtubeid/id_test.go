package youtubeid_test

import (
	"testing"

	"github.com/xynova/library-intake/internal/youtubeid"
)

func TestFromURL(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in   string
		want string
	}{
		{"nZYJdwM-_nI", "nZYJdwM-_nI"},
		{"https://www.youtube.com/watch?v=nZYJdwM-_nI", "nZYJdwM-_nI"},
		{"https://youtu.be/nZYJdwM-_nI", "nZYJdwM-_nI"},
		{"https://www.youtube.com/embed/nZYJdwM-_nI", "nZYJdwM-_nI"},
		{"https://www.youtube.com/shorts/nZYJdwM-_nI", "nZYJdwM-_nI"},
	}
	for _, tc := range cases {
		got, err := youtubeid.FromURL(tc.in)
		if err != nil {
			t.Fatalf("%q: %v", tc.in, err)
		}
		if got != tc.want {
			t.Fatalf("%q: got %q want %q", tc.in, got, tc.want)
		}
	}
}

func TestFromURLInvalid(t *testing.T) {
	t.Parallel()
	if _, err := youtubeid.FromURL("https://example.com/watch?v=nope"); err == nil {
		t.Fatal("expected error")
	}
}
