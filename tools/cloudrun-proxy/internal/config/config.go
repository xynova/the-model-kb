package config

import (
	"fmt"
	"math"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ConfigFile string
	Server     ServerSettings
	Auth       AuthSettings
	Defaults   ModelDefaults
	Models     map[string]ModelSpec
}

// ResolvedModel is a fully merged model route ready for proxying.
type ResolvedModel struct {
	ID                   string
	Audience             string
	URL                  string
	ChatPath             string
	RequestTimeout       time.Duration
	ImpersonateSA       string
	UpstreamModel        string
	EnableThinking       bool
	MinMaxTokens         int
	MinMaxTokensThinking int
}

type ServerSettings struct {
	ListenAddr         string
	Debug              bool
	LocalSecret        string
	MaxConcurrent      int
	RejectWhenBusy     bool
	RateLimitCooldown  time.Duration
	rateLimitCooldownRaw string
}

type AuthSettings struct {
	ImpersonateServiceAccount string
}

type ModelDefaults struct {
	ChatPath             string
	RequestTimeout       time.Duration
	requestTimeoutRaw    string
	UpstreamModel        string
	EnableThinking       bool
	MinMaxTokens         int
	MinMaxTokensThinking int
}

type ModelSpec struct {
	Audience                  string
	URL                       string
	ImpersonateServiceAccount string
	ChatPath                  string
	RequestTimeout            time.Duration
	requestTimeoutRaw         string
	UpstreamModel             string
	EnableThinking            *bool
	MinMaxTokens              int
	MinMaxTokensThinking      int
}

type fileConfig struct {
	Server   serverConfig            `yaml:"server"`
	Auth     authConfig              `yaml:"auth"`
	Defaults modelDefaultsConfig     `yaml:"defaults"`
	Models   map[string]modelYAML    `yaml:"models"`
}

type serverConfig struct {
	ListenAddr          string `yaml:"listen_addr"`
	Debug               *bool  `yaml:"debug"`
	LocalSecret         string `yaml:"local_secret"`
	MaxConcurrent       int    `yaml:"max_concurrent"`
	RejectWhenBusy      *bool  `yaml:"reject_when_busy"`
	RateLimitCooldown   string `yaml:"rate_limit_cooldown"`
}

type authConfig struct {
	ImpersonateServiceAccount string `yaml:"impersonate_service_account"`
}

type modelDefaultsConfig struct {
	ChatPath             string `yaml:"chat_path"`
	RequestTimeout       string `yaml:"request_timeout"`
	UpstreamModel        string `yaml:"upstream_model"`
	EnableThinking       *bool  `yaml:"enable_thinking"`
	MinMaxTokens         int    `yaml:"min_max_tokens"`
	MinMaxTokensThinking int    `yaml:"min_max_tokens_thinking"`
}

type modelYAML struct {
	Audience                  string `yaml:"audience"`
	URL                       string `yaml:"url"`
	ImpersonateServiceAccount string `yaml:"impersonate_service_account"`
	ChatPath                  string `yaml:"chat_path"`
	RequestTimeout            string `yaml:"request_timeout"`
	UpstreamModel             string `yaml:"upstream_model"`
	EnableThinking            *bool  `yaml:"enable_thinking"`
	MinMaxTokens              int    `yaml:"min_max_tokens"`
	MinMaxTokensThinking      int    `yaml:"min_max_tokens_thinking"`
}

func Load() (Config, error) {
	cfg := defaultConfig()

	path := configFilePath()
	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			if !os.IsNotExist(err) {
				return Config{}, fmt.Errorf("read config file %s: %w", path, err)
			}
		} else {
			var file fileConfig
			if err := yaml.Unmarshal(data, &file); err != nil {
				return Config{}, fmt.Errorf("parse config file %s: %w", path, err)
			}
			file.mergeInto(&cfg)
			cfg.ConfigFile = path
		}
	}

	if err := parseDurations(&cfg); err != nil {
		return Config{}, err
	}

	applyEnv(&cfg)

	if err := parseDurations(&cfg); err != nil {
		return Config{}, err
	}

	if err := cfg.validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c Config) ModelIDs() []string {
	ids := make([]string, 0, len(c.Models))
	for id := range c.Models {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

func (c Config) Resolve(id string) (ResolvedModel, bool) {
	spec, ok := c.Models[id]
	if !ok {
		return ResolvedModel{}, false
	}
	return spec.resolve(id, c.Defaults, c.Auth), true
}

func (s ModelSpec) resolve(id string, def ModelDefaults, auth AuthSettings) ResolvedModel {
	m := ResolvedModel{
		ID:                   id,
		Audience:             strings.TrimRight(strings.TrimSpace(s.Audience), "/"),
		URL:                  strings.TrimRight(strings.TrimSpace(s.URL), "/"),
		ChatPath:             def.ChatPath,
		RequestTimeout:       def.RequestTimeout,
		ImpersonateSA:        auth.ImpersonateServiceAccount,
		UpstreamModel:        def.UpstreamModel,
		EnableThinking:       def.EnableThinking,
		MinMaxTokens:         def.MinMaxTokens,
		MinMaxTokensThinking: def.MinMaxTokensThinking,
	}

	if v := strings.TrimSpace(s.ChatPath); v != "" {
		m.ChatPath = v
	}
	if s.RequestTimeout > 0 {
		m.RequestTimeout = s.RequestTimeout
	}
	if v := strings.TrimSpace(s.UpstreamModel); v != "" {
		m.UpstreamModel = v
	}
	if s.EnableThinking != nil {
		m.EnableThinking = *s.EnableThinking
	}
	if s.MinMaxTokens > 0 {
		m.MinMaxTokens = s.MinMaxTokens
	}
	if s.MinMaxTokensThinking > 0 {
		m.MinMaxTokensThinking = s.MinMaxTokensThinking
	}
	if v := strings.TrimSpace(s.ImpersonateServiceAccount); v != "" {
		m.ImpersonateSA = v
	}
	return m
}

func (f fileConfig) mergeInto(cfg *Config) {
	if v := strings.TrimSpace(f.Server.ListenAddr); v != "" {
		cfg.Server.ListenAddr = v
	}
	if f.Server.Debug != nil {
		cfg.Server.Debug = *f.Server.Debug
	}
	if v := strings.TrimSpace(f.Server.LocalSecret); v != "" {
		cfg.Server.LocalSecret = v
	}
	if f.Server.MaxConcurrent > 0 {
		cfg.Server.MaxConcurrent = f.Server.MaxConcurrent
	}
	if f.Server.RejectWhenBusy != nil {
		cfg.Server.RejectWhenBusy = *f.Server.RejectWhenBusy
	}
	if v := strings.TrimSpace(f.Server.RateLimitCooldown); v != "" {
		cfg.Server.rateLimitCooldownRaw = v
	}

	if v := strings.TrimSpace(f.Auth.ImpersonateServiceAccount); v != "" {
		cfg.Auth.ImpersonateServiceAccount = v
	}

	if v := strings.TrimSpace(f.Defaults.ChatPath); v != "" {
		cfg.Defaults.ChatPath = v
	}
	if v := strings.TrimSpace(f.Defaults.RequestTimeout); v != "" {
		cfg.Defaults.requestTimeoutRaw = v
	}
	if v := strings.TrimSpace(f.Defaults.UpstreamModel); v != "" {
		cfg.Defaults.UpstreamModel = v
	}
	if f.Defaults.EnableThinking != nil {
		cfg.Defaults.EnableThinking = *f.Defaults.EnableThinking
	}
	if f.Defaults.MinMaxTokens > 0 {
		cfg.Defaults.MinMaxTokens = f.Defaults.MinMaxTokens
	}
	if f.Defaults.MinMaxTokensThinking > 0 {
		cfg.Defaults.MinMaxTokensThinking = f.Defaults.MinMaxTokensThinking
	}

	if len(f.Models) > 0 {
		cfg.Models = make(map[string]ModelSpec, len(f.Models))
		for id, raw := range f.Models {
			id = strings.TrimSpace(id)
			if id == "" {
				continue
			}
			cfg.Models[id] = raw.toSpec()
		}
	}
}

func (m modelYAML) toSpec() ModelSpec {
	return ModelSpec{
		Audience:                  m.Audience,
		URL:                       m.URL,
		ImpersonateServiceAccount: m.ImpersonateServiceAccount,
		ChatPath:                  m.ChatPath,
		requestTimeoutRaw:         m.RequestTimeout,
		UpstreamModel:             m.UpstreamModel,
		EnableThinking:            m.EnableThinking,
		MinMaxTokens:              m.MinMaxTokens,
		MinMaxTokensThinking:      m.MinMaxTokensThinking,
	}
}

func defaultConfig() Config {
	return Config{
		Server: ServerSettings{
			ListenAddr:        "127.0.0.1:4318",
			MaxConcurrent:     1,
			RejectWhenBusy:    true,
			rateLimitCooldownRaw: "60s",
			RateLimitCooldown: 60 * time.Second,
		},
		Defaults: ModelDefaults{
			ChatPath:             "/v1/chat/completions",
			requestTimeoutRaw:    "900s",
			RequestTimeout:       900 * time.Second,
			UpstreamModel:        "gemma-4-12b",
			MinMaxTokens:         512,
			MinMaxTokensThinking: 2048,
		},
		Models: map[string]ModelSpec{},
	}
}

func parseDurations(cfg *Config) error {
	if err := parseDurationField(&cfg.Defaults.requestTimeoutRaw, &cfg.Defaults.RequestTimeout, "defaults.request_timeout"); err != nil {
		return err
	}
	if err := parseDurationField(&cfg.Server.rateLimitCooldownRaw, &cfg.Server.RateLimitCooldown, "server.rate_limit_cooldown"); err != nil {
		return err
	}
	for id, spec := range cfg.Models {
		if spec.requestTimeoutRaw == "" {
			spec.RequestTimeout = cfg.Defaults.RequestTimeout
			cfg.Models[id] = spec
			continue
		}
		if err := parseDurationField(&spec.requestTimeoutRaw, &spec.RequestTimeout, fmt.Sprintf("models.%s.request_timeout", id)); err != nil {
			return err
		}
		cfg.Models[id] = spec
	}
	return nil
}

func parseDurationField(raw *string, dest *time.Duration, field string) error {
	s := strings.TrimSpace(*raw)
	if s == "" {
		return nil
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("%s: %w", field, err)
	}
	*dest = d
	return nil
}

func configFilePath() string {
	if path := strings.TrimSpace(os.Getenv("CONFIG_FILE")); path != "" {
		return path
	}
	return "config.yaml"
}

func applyEnv(cfg *Config) {
	if v := strings.TrimSpace(os.Getenv("LISTEN_ADDR")); v != "" {
		cfg.Server.ListenAddr = v
	}
	if v, ok := os.LookupEnv("LOCAL_SECRET"); ok {
		cfg.Server.LocalSecret = v
	}
	if v, ok := envBool("DEBUG"); ok {
		cfg.Server.Debug = v
	}
	if v, ok := envInt("GEMMA_MAX_CONCURRENT"); ok {
		cfg.Server.MaxConcurrent = v
	}
	if v := strings.TrimSpace(os.Getenv("IMPERSONATE_SERVICE_ACCOUNT")); v != "" {
		cfg.Auth.ImpersonateServiceAccount = v
	}
}

func (c Config) validate() error {
	if len(c.Models) == 0 {
		return fmt.Errorf("models: at least one model route is required")
	}
	if err := validateLoopbackOnly(c.Server.ListenAddr); err != nil {
		return err
	}
	if c.Server.MaxConcurrent <= 0 {
		return fmt.Errorf("server.max_concurrent must be positive")
	}
	if !strings.HasPrefix(c.Defaults.ChatPath, "/") {
		return fmt.Errorf("defaults.chat_path must start with /")
	}
	if c.Defaults.RequestTimeout <= 0 {
		return fmt.Errorf("defaults.request_timeout must be a positive duration")
	}

	for id, spec := range c.Models {
		resolved := spec.resolve(id, c.Defaults, c.Auth)
		if resolved.Audience == "" {
			return fmt.Errorf("models.%s.audience is required", id)
		}
		if resolved.URL == "" {
			return fmt.Errorf("models.%s.url is required", id)
		}
		if !strings.HasPrefix(resolved.ChatPath, "/") {
			return fmt.Errorf("models.%s.chat_path must start with /", id)
		}
		if resolved.RequestTimeout <= 0 {
			return fmt.Errorf("models.%s.request_timeout must be a positive duration", id)
		}
	}
	return nil
}

func validateLoopbackOnly(listenAddr string) error {
	host, _, err := net.SplitHostPort(listenAddr)
	if err != nil {
		return fmt.Errorf("server.listen_addr: %w", err)
	}
	if host != "127.0.0.1" && host != "localhost" {
		return fmt.Errorf("server.listen_addr must bind to loopback (127.0.0.1), got %q", host)
	}
	return nil
}

func envBool(key string) (bool, bool) {
	raw, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(raw) == "" {
		return false, false
	}
	v, err := strconv.ParseBool(raw)
	if err != nil {
		return false, false
	}
	return v, true
}

func envInt(key string) (int, bool) {
	raw, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(raw) == "" {
		return 0, false
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n <= 0 || n > math.MaxInt32 {
		return 0, false
	}
	return n, true
}
