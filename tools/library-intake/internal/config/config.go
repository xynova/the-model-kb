package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultBaseURL       = "http://127.0.0.1:1320"
	DefaultChatModel     = "cf_local/@cf/zai-org/glm-4.7-flash"
	DefaultVisionModel   = "cf_local/@cf/google/gemma-4-26b-a4b-it"
	DefaultFrameEvery    = 8
	DefaultVisionBatch   = 4
	DefaultWorkRoot      = "tmp/library-intake"
	DefaultChatTimeout   = 180 * time.Second
	DefaultVisionTimeout = 300 * time.Second
	// DefaultYtDlpExtractorArgs prefers clients that still serve progressive
	// HTTPS URLs. Default android_vr / SABR-only sessions often 403.
	DefaultYtDlpExtractorArgs = "youtube:player_client=android,web"
)

// Config holds construction inputs for the library-intake CLI services.
type Config struct {
	BaseURL            string
	ChatModel          string
	VisionModel        string
	WorkRoot           string
	FrameIntervalS     int
	VisionBatch        int
	ChatTimeout        time.Duration
	VisionTimeout      time.Duration
	FabricBin          string
	YtDlpBin           string
	FFmpegBin          string
	YtDlpExtractorArgs string
}

// Load reads environment overrides into a Config.
func Load() Config {
	cfg := Config{
		BaseURL:            envOr("POLYPUS_BASE_URL", DefaultBaseURL),
		ChatModel:          envOr("LIBRARY_INTAKE_POLYPUS_MODEL", DefaultChatModel),
		VisionModel:        envOr("LIBRARY_INTAKE_VISION_MODEL", DefaultVisionModel),
		WorkRoot:           envOr("LIBRARY_INTAKE_WORK_ROOT", DefaultWorkRoot),
		FrameIntervalS:     envInt("LIBRARY_INTAKE_FRAME_INTERVAL", DefaultFrameEvery),
		VisionBatch:        envInt("LIBRARY_INTAKE_VISION_BATCH", DefaultVisionBatch),
		ChatTimeout:        DefaultChatTimeout,
		VisionTimeout:      DefaultVisionTimeout,
		FabricBin:          envOr("LIBRARY_INTAKE_FABRIC", "fabric"),
		YtDlpBin:           envOr("LIBRARY_INTAKE_YTDLP", "yt-dlp"),
		FFmpegBin:          envOr("LIBRARY_INTAKE_FFMPEG", "ffmpeg"),
		YtDlpExtractorArgs: envOr("LIBRARY_INTAKE_YTDLP_EXTRACTOR_ARGS", DefaultYtDlpExtractorArgs),
	}
	return cfg
}

func envOr(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return fallback
	}
	return n
}
