package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const adcTypeAuthorizedUser = "authorized_user"

type adcMeta struct {
	Type string `json:"type"`
}

func adcPath() (string, error) {
	if path := strings.TrimSpace(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")); path != "" {
		return path, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("home dir: %w", err)
	}
	return filepath.Join(home, ".config", "gcloud", "application_default_credentials.json"), nil
}

func readADCMeta() (adcMeta, string, error) {
	path, err := adcPath()
	if err != nil {
		return adcMeta{}, "", err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return adcMeta{}, path, fmt.Errorf("read ADC file %s: %w", path, err)
	}

	var meta adcMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		return adcMeta{}, path, fmt.Errorf("parse ADC file: %w", err)
	}
	return meta, path, nil
}

func isPersonalADC() bool {
	meta, _, err := readADCMeta()
	return err == nil && meta.Type == adcTypeAuthorizedUser
}
