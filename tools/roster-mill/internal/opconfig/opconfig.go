package opconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	derrors "github.com/xynova/roster-mill/internal/errors"
	"github.com/zalando/go-keyring"
	"gopkg.in/yaml.v3"
)

const (
	AppDir      = "the-model-kb"
	OverrideEnv = "THE_MODEL_KB_CONFIG"
	SecretAA    = "AA_API_KEY"
	Placeholder = "${AA_API_KEY}"
)

// ExampleYAML is the default operator config body.
const ExampleYAML = `# Live operator config. Mode 0600. No literal secrets.
secrets:
  - env: AA_API_KEY
    required: true
aa:
  api_key: "${AA_API_KEY}"
`

// Keyring is a testable secret store.
type Keyring interface {
	Get(service, account string) (string, error)
	Set(service, account, value string) error
}

// OSKeyring is the zalando/go-keyring backend.
type OSKeyring struct{}

func (OSKeyring) Get(service, account string) (string, error) {
	return keyring.Get(service, account)
}

func (OSKeyring) Set(service, account, value string) error {
	return keyring.Set(service, account, value)
}

// MemKeyring is an in-memory store for tests.
type MemKeyring struct {
	Store map[string]string
}

func (m *MemKeyring) key(service, account string) string {
	return service + "\x00" + account
}

func (m *MemKeyring) Get(service, account string) (string, error) {
	if m.Store == nil {
		return "", keyring.ErrNotFound
	}
	v, ok := m.Store[m.key(service, account)]
	if !ok {
		return "", keyring.ErrNotFound
	}
	return v, nil
}

func (m *MemKeyring) Set(service, account, value string) error {
	if m.Store == nil {
		m.Store = map[string]string{}
	}
	m.Store[m.key(service, account)] = value
	return nil
}

func getSecret(ring Keyring, service, account string) string {
	if ring == nil {
		return ""
	}
	v, err := ring.Get(service, account)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(v)
}

// UserConfigDir is ~/.config/the-model-kb or $XDG_CONFIG_HOME/the-model-kb.
func UserConfigDir(xdgHome string) string {
	if xdgHome != "" {
		return filepath.Join(xdgHome, AppDir)
	}
	if xdg := strings.TrimSpace(os.Getenv("XDG_CONFIG_HOME")); xdg != "" {
		return filepath.Join(xdg, AppDir)
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", AppDir)
}

// ConfigCandidates is override, then user config, then repo site/config.yaml.
func ConfigCandidates(override, xdgHome, repoRoot string) []string {
	var paths []string
	raw := override
	if raw == "" {
		raw = strings.TrimSpace(os.Getenv(OverrideEnv))
	}
	if raw != "" {
		paths = append(paths, raw)
	}
	paths = append(paths, filepath.Join(UserConfigDir(xdgHome), "config.yaml"))
	if repoRoot != "" {
		paths = append(paths, filepath.Join(repoRoot, "site", "config.yaml"))
	}
	return paths
}

// LoadDotenvFile fills missing keys from a KEY=VALUE file.
func LoadDotenvFile(path string, environ map[string]string) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(raw), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}
		key, value, _ := strings.Cut(line, "=")
		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		if key != "" {
			if _, exists := environ[key]; !exists {
				environ[key] = value
			}
		}
	}
}

// LoadDotenvFiles loads repo and site .env into environ.
func LoadDotenvFiles(repoRoot string, environ map[string]string) {
	LoadDotenvFile(filepath.Join(repoRoot, ".env"), environ)
	LoadDotenvFile(filepath.Join(repoRoot, "site", ".env"), environ)
}

func declaredSecrets(data map[string]any) []struct {
	env      string
	required bool
} {
	rows, _ := data["secrets"].([]any)
	var out []struct {
		env      string
		required bool
	}
	for _, row := range rows {
		if s, ok := row.(string); ok {
			out = append(out, struct {
				env      string
				required bool
			}{env: s})
			continue
		}
		if m, ok := row.(map[string]any); ok {
			env, _ := m["env"].(string)
			if env != "" {
				req, _ := m["required"].(bool)
				out = append(out, struct {
					env      string
					required bool
				}{env: env, required: req})
			}
		}
	}
	return out
}

func expandPlaceholders(value string, environ map[string]string) string {
	out := value
	for key, val := range environ {
		out = strings.ReplaceAll(out, "${"+key+"}", val)
	}
	return out
}

func getenv(environ map[string]string, key string) string {
	if environ != nil {
		return strings.TrimSpace(environ[key])
	}
	return strings.TrimSpace(os.Getenv(key))
}

func setenv(environ map[string]string, key, value string) {
	if environ != nil {
		environ[key] = value
		return
	}
	_ = os.Setenv(key, value)
}

// ResolveOptions is AA key lookup input.
type ResolveOptions struct {
	Environ  map[string]string
	Ring     Keyring
	Override string
	XDGHome  string
	RepoRoot string
}

// ResolveAAAPIKey follows env, then keyring + config.yaml placeholders.
func ResolveAAAPIKey(opts ResolveOptions) (string, error) {
	fromEnv := getenv(opts.Environ, SecretAA)
	if fromEnv != "" {
		return fromEnv, nil
	}
	ring := opts.Ring
	if ring == nil {
		ring = OSKeyring{}
	}
	for _, path := range ConfigCandidates(opts.Override, opts.XDGHome, opts.RepoRoot) {
		raw, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var data map[string]any
		if err := yaml.Unmarshal(raw, &data); err != nil || data == nil {
			continue
		}
		names := declaredSecrets(data)
		if len(names) == 0 {
			continue
		}
		required := false
		for _, n := range names {
			if n.env != SecretAA {
				continue
			}
			required = required || n.required
			if stored := getSecret(ring, AppDir, n.env); stored != "" {
				setenv(opts.Environ, n.env, stored)
			}
		}
		aa, _ := data["aa"].(map[string]any)
		rawKey := ""
		if aa != nil {
			rawKey = fmt.Sprint(aa["api_key"])
			if aa["api_key"] == nil {
				rawKey = ""
			}
		}
		envCopy := map[string]string{}
		if opts.Environ != nil {
			for k, v := range opts.Environ {
				envCopy[k] = v
			}
		} else {
			for _, e := range os.Environ() {
				k, v, _ := strings.Cut(e, "=")
				envCopy[k] = v
			}
		}
		expanded := strings.TrimSpace(expandPlaceholders(rawKey, envCopy))
		if expanded != "" && !strings.Contains(expanded, Placeholder) {
			return expanded, nil
		}
		if required && getenv(opts.Environ, SecretAA) == "" {
			return "", derrors.New(derrors.CodeInvalidArgument, "opconfig",
				SecretAA+" required in "+path+" but env and keyring are empty")
		}
		break
	}
	return getenv(opts.Environ, SecretAA), nil
}

// InitUserConfig writes config.yaml.example and optionally live config.yaml.
func InitUserConfig(xdgHome string, force bool, exampleSrc string) (live, example string, err error) {
	destDir := UserConfigDir(xdgHome)
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", "", err
	}
	live = filepath.Join(destDir, "config.yaml")
	example = filepath.Join(destDir, "config.yaml.example")
	body := ExampleYAML
	if exampleSrc != "" {
		if raw, err := os.ReadFile(exampleSrc); err == nil {
			body = string(raw)
		}
	}
	if err := os.WriteFile(example, []byte(body), 0o644); err != nil {
		return "", "", err
	}
	if force {
		if err := os.WriteFile(live, []byte(body), 0o600); err != nil {
			return "", "", err
		}
		return live, example, nil
	}
	if _, err := os.Stat(live); err == nil {
		return live, example, nil
	}
	if err := os.WriteFile(live, []byte(body), 0o600); err != nil {
		return "", "", err
	}
	return live, example, nil
}

// StoreAAKey writes AA_API_KEY into the OS keyring.
func StoreAAKey(value string, ring Keyring) error {
	secret := strings.TrimSpace(value)
	if secret == "" {
		return derrors.New(derrors.CodeInvalidArgument, "opconfig", SecretAA+" is empty")
	}
	if ring == nil {
		ring = OSKeyring{}
	}
	if err := ring.Set(AppDir, SecretAA, secret); err != nil {
		return derrors.Wrap(err, derrors.CodeFailed, "opconfig", "keyring write")
	}
	return nil
}

// EnvironFromOS copies os.Environ into a map.
func EnvironFromOS() map[string]string {
	out := map[string]string{}
	for _, e := range os.Environ() {
		k, v, _ := strings.Cut(e, "=")
		out[k] = v
	}
	return out
}
