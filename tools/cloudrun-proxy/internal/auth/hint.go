package auth

import "strings"

// ErrorHint returns a short remediation hint for common auth failures.
func ErrorHint(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()
	switch {
	case strings.Contains(msg, "invalid_grant"),
		strings.Contains(msg, "invalid_rapt"),
		strings.Contains(msg, "reauth related error"):
		return "ADC session expired — run: gcloud auth application-default login"
	case strings.Contains(msg, "gcloud auth print-identity-token"):
		return "gcloud user session expired — run: gcloud auth login && gcloud auth application-default login"
	case strings.Contains(msg, "Permission") && strings.Contains(msg, "denied"),
		strings.Contains(msg, "403"):
		return "check IAM: serviceAccountTokenCreator on invoker SA, run.invoker on Cloud Run service"
	default:
		return ""
	}
}
