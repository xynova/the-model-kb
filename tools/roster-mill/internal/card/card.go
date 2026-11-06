package card

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Keep in sync with roster/taxonomy.md.
var (
	AllowedKinds = map[string]struct{}{
		"sovereign": {}, "provider": {}, "service": {},
	}
	AllowedCategories = map[string]struct{}{
		"cheap": {}, "mid": {}, "spend": {}, "unpriced": {}, "classify": {}, "tts": {},
	}
	AllowedStatus = map[string]struct{}{
		"reach-for": {}, "watch": {}, "retired": {},
	}
	AllowedJobs = map[string]struct{}{
		"chat": {}, "reasoning": {}, "coding": {}, "vision": {},
		"long-context": {}, "cheap": {}, "high-throughput": {},
		"classify": {}, "tts": {},
	}
	AllowedRuntimes = map[string]struct{}{
		"cloudflare": {}, "deepinfra": {}, "openrouter": {},
		"google": {}, "xinference": {}, "local": {},
	}
	AllowedCapacity = map[string]struct{}{
		"hobby": {}, "industrial": {}, "local": {}, "unknown": {},
	}
	AllowedHostIDs = map[string]struct{}{
		"cloudflare": {}, "deepinfra": {}, "openrouter": {},
		"google": {}, "xinference": {}, "cursor": {}, "local": {},
	}
	SkipNames = map[string]struct{}{
		"README.md": {}, "INDEX.md": {}, "taxonomy.md": {},
		"_template.md": {}, "OBJECTIVE.md": {}, "catalog.json": {},
	}
)

const (
	CheapMax = 3.0
	MidMax   = 15.0
)

var cardRequired = []string{
	"name", "title", "kind", "category", "jobs", "tags",
	"status", "lab", "source", "refreshed", "hosts",
}

// File is a roster card with YAML frontmatter preserved as a node tree.
type File struct {
	Front *yaml.Node
	Body  string
}

// Split parses --- frontmatter. On failure Front is nil and Body is the original text.
func Split(text string) *File {
	if !strings.HasPrefix(text, "---") {
		return &File{Body: text}
	}
	end := strings.Index(text[3:], "\n---")
	if end < 0 {
		return &File{Body: text}
	}
	end += 3
	fmRaw := text[3:end]
	body := text[end+4:]
	if strings.HasPrefix(body, "\n") {
		body = body[1:]
	}
	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(fmRaw), &doc); err != nil {
		return &File{Body: text}
	}
	if Mapping(&doc) == nil {
		return &File{Body: text}
	}
	return &File{Front: &doc, Body: body}
}

// Mapping returns the root mapping node of a YAML document.
func Mapping(doc *yaml.Node) *yaml.Node {
	if doc == nil {
		return nil
	}
	if doc.Kind == yaml.DocumentNode && len(doc.Content) > 0 {
		if doc.Content[0] != nil && doc.Content[0].Kind == yaml.MappingNode {
			return doc.Content[0]
		}
		return nil
	}
	if doc.Kind == yaml.MappingNode {
		return doc
	}
	return nil
}

// Lookup returns the value node for key, or nil.
func Lookup(m *yaml.Node, key string) *yaml.Node {
	if m == nil {
		return nil
	}
	for i := 0; i+1 < len(m.Content); i += 2 {
		if m.Content[i].Value == key {
			return m.Content[i+1]
		}
	}
	return nil
}

// Set replaces or appends key in mapping m.
func Set(m *yaml.Node, key string, val *yaml.Node) {
	if m == nil || val == nil {
		return
	}
	for i := 0; i+1 < len(m.Content); i += 2 {
		if m.Content[i].Value == key {
			m.Content[i+1] = val
			return
		}
	}
	m.Content = append(m.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: key},
		val,
	)
}

// CloneNode deep-copies n.
func CloneNode(n *yaml.Node) *yaml.Node {
	if n == nil {
		return nil
	}
	c := *n
	if n.Content != nil {
		c.Content = make([]*yaml.Node, len(n.Content))
		for i, ch := range n.Content {
			c.Content[i] = CloneNode(ch)
		}
	}
	return &c
}

// AsMap decodes frontmatter into a generic map.
func (f *File) AsMap() map[string]any {
	m := Mapping(f.Front)
	if m == nil {
		return nil
	}
	var out map[string]any
	if err := m.Decode(&out); err != nil {
		return nil
	}
	return out
}

// DumpSource encodes frontmatter without Hugo remapping.
func DumpSource(front *yaml.Node) (string, error) {
	return encodeMapping(Mapping(front))
}

// DumpHugo remaps kind to roster_kind and quotes numeric aa_intelligence.
func DumpHugo(front *yaml.Node) (string, error) {
	cloned := CloneNode(front)
	m := Mapping(cloned)
	if m == nil {
		return "", fmt.Errorf("frontmatter is not a mapping")
	}
	for i := 0; i+1 < len(m.Content); i += 2 {
		if m.Content[i].Value == "kind" {
			m.Content[i].Value = "roster_kind"
		}
	}
	if v := Lookup(m, "aa_intelligence"); v != nil {
		var raw any
		if err := v.Decode(&raw); err == nil {
			if s := formatQuotedScore(raw); s != "" {
				v.Kind = yaml.ScalarNode
				v.Tag = "!!str"
				v.Value = s
				v.Style = yaml.DoubleQuotedStyle
			}
		}
	}
	return encodeMapping(m)
}

func encodeMapping(m *yaml.Node) (string, error) {
	if m == nil {
		return "", fmt.Errorf("nil mapping")
	}
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(m); err != nil {
		_ = enc.Close()
		return "", err
	}
	if err := enc.Close(); err != nil {
		return "", err
	}
	return strings.TrimRight(buf.String(), "\n"), nil
}

func formatQuotedScore(raw any) string {
	switch t := raw.(type) {
	case bool:
		return ""
	case int:
		return fmt.Sprintf("%d", t)
	case int64:
		return fmt.Sprintf("%d", t)
	case uint64:
		return fmt.Sprintf("%d", t)
	case float32:
		return trimFloat(float64(t))
	case float64:
		return trimFloat(t)
	default:
		return ""
	}
}

func trimFloat(f float64) string {
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.15g", f), "0"), ".")
}

// ScalarString is a YAML string scalar.
func ScalarString(v string) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: v}
}

// ScalarQuoted is a double-quoted YAML string.
func ScalarQuoted(v string) *yaml.Node {
	n := ScalarString(v)
	n.Style = yaml.DoubleQuotedStyle
	return n
}

// WriteCard writes --- frontmatter --- body using source dump (no kind remap).
func WriteCard(path string, front *yaml.Node, body string) error {
	dumped, err := DumpSource(front)
	if err != nil {
		return err
	}
	sep := "\n"
	if body != "" && !strings.HasPrefix(body, "\n") {
		sep = "\n\n"
	}
	text := "---\n" + dumped + "\n---" + sep + body
	return os.WriteFile(path, []byte(text), 0o644)
}

// FindByName walks cardsRoot for a markdown card whose name matches.
func FindByName(cardsRoot, name string) (string, *File, error) {
	var found string
	var parsed *File
	err := filepath.WalkDir(cardsRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		base := d.Name()
		if !strings.HasSuffix(base, ".md") || strings.HasPrefix(base, "_") {
			return nil
		}
		raw, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		f := Split(string(raw))
		fm := f.AsMap()
		if fm == nil {
			return nil
		}
		n, _ := fm["name"].(string)
		if n == name {
			found = path
			parsed = f
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		return "", nil, err
	}
	if found == "" {
		return "", nil, nil
	}
	return found, parsed, nil
}

// HostMapping finds hosts[i] whose id equals hostID.
func HostMapping(fm map[string]any, hostsNode *yaml.Node, hostID string) (map[string]any, *yaml.Node) {
	hosts, _ := fm["hosts"].([]any)
	if hostsNode == nil || hostsNode.Kind != yaml.SequenceNode {
		return nil, nil
	}
	for i, h := range hosts {
		hm, ok := h.(map[string]any)
		if !ok {
			continue
		}
		id, _ := hm["id"].(string)
		if id == hostID && i < len(hostsNode.Content) {
			return hm, hostsNode.Content[i]
		}
	}
	return nil, nil
}
