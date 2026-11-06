package opconfig

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnvWinsOverKeyring(t *testing.T) {
	ring := &MemKeyring{}
	_ = ring.Set(AppDir, SecretAA, "from-keyring")
	got, err := ResolveAAAPIKey(ResolveOptions{
		Environ: map[string]string{"AA_API_KEY": "from-env"},
		Ring:    ring,
	})
	if err != nil || got != "from-env" {
		t.Fatalf("got=%q err=%v", got, err)
	}
}

func TestOverrideConfigThenKeyring(t *testing.T) {
	ring := &MemKeyring{}
	_ = ring.Set(AppDir, SecretAA, "from-keyring")
	cfg := filepath.Join(t.TempDir(), "config.yaml")
	_ = os.WriteFile(cfg, []byte(ExampleYAML), 0o644)
	got, err := ResolveAAAPIKey(ResolveOptions{
		Environ:  map[string]string{},
		Ring:     ring,
		Override: cfg,
	})
	if err != nil || got != "from-keyring" {
		t.Fatalf("got=%q err=%v", got, err)
	}
}

func TestRequiredEmptyFailsClosed(t *testing.T) {
	cfg := filepath.Join(t.TempDir(), "config.yaml")
	_ = os.WriteFile(cfg, []byte(ExampleYAML), 0o644)
	_, err := ResolveAAAPIKey(ResolveOptions{
		Environ:  map[string]string{},
		Ring:     &MemKeyring{},
		Override: cfg,
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDotenvFillsMissingOnly(t *testing.T) {
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, ".env"), []byte("AA_API_KEY=from-dotenv\n"), 0o644)
	env := map[string]string{"AA_API_KEY": "already"}
	LoadDotenvFiles(root, env)
	if env["AA_API_KEY"] != "already" {
		t.Fatal(env)
	}
	env2 := map[string]string{}
	LoadDotenvFiles(root, env2)
	if env2["AA_API_KEY"] != "from-dotenv" {
		t.Fatal(env2)
	}
}

func TestInitWrites0600AndDoesNotOverwrite(t *testing.T) {
	tmp := t.TempDir()
	live, example, err := InitUserConfig(tmp, false, "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(example); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(live)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("mode=%o", info.Mode().Perm())
	}
	_ = os.WriteFile(live, []byte("keep\n"), 0o644)
	_, _, err = InitUserConfig(tmp, false, "")
	if err != nil {
		t.Fatal(err)
	}
	body, _ := os.ReadFile(live)
	if string(body) != "keep\n" {
		t.Fatalf("%q", body)
	}
}

func TestStoreAAKey(t *testing.T) {
	ring := &MemKeyring{}
	if err := StoreAAKey("secret", ring); err != nil {
		t.Fatal(err)
	}
	got, err := ring.Get(AppDir, SecretAA)
	if err != nil || got != "secret" {
		t.Fatalf("%q %v", got, err)
	}
	if err := StoreAAKey("  ", ring); err == nil {
		t.Fatal("expected empty error")
	}
}
