package youtubeid

import (
	"net/url"
	"regexp"
	"strings"

	derrors "github.com/xynova/library-intake/internal/errors"
)

var idRe = regexp.MustCompile(`^[A-Za-z0-9_-]{11}$`)

// FromURL extracts an 11-character YouTube video id from common URL shapes.
func FromURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", derrors.New(derrors.CodeInvalidArgument, "youtubeid.FromURL", "url is empty")
	}
	if idRe.MatchString(raw) {
		return raw, nil
	}
	u, err := url.Parse(raw)
	if err != nil {
		return "", derrors.Wrap(err, derrors.CodeInvalidArgument, "youtubeid.FromURL", "parse url")
	}
	host := strings.ToLower(u.Host)
	switch {
	case strings.Contains(host, "youtu.be"):
		id := strings.Trim(strings.TrimPrefix(u.Path, "/"), "/")
		if i := strings.IndexByte(id, '/'); i >= 0 {
			id = id[:i]
		}
		if idRe.MatchString(id) {
			return id, nil
		}
	case strings.Contains(host, "youtube.com"), strings.Contains(host, "youtube-nocookie.com"):
		if v := u.Query().Get("v"); idRe.MatchString(v) {
			return v, nil
		}
		parts := strings.Split(strings.Trim(u.Path, "/"), "/")
		if len(parts) >= 2 && (parts[0] == "embed" || parts[0] == "shorts" || parts[0] == "live") {
			if idRe.MatchString(parts[1]) {
				return parts[1], nil
			}
		}
	}
	return "", derrors.New(derrors.CodeInvalidArgument, "youtubeid.FromURL", "could not extract video id").
		With("url", raw)
}
