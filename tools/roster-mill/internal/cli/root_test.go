package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestBareInvokePrintsAgentGuide(t *testing.T) {
	root := newRootCmd()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "Documentation for agents") {
		t.Fatalf("%s", out)
	}
	if !strings.Contains(out, "generate") {
		t.Fatalf("%s", out)
	}
}

func TestVersion(t *testing.T) {
	root := newRootCmd()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"version"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(buf.String()) != "dev" {
		t.Fatalf("%q", buf.String())
	}
}
