package store

import (
	"errors"
	"fmt"
	"strings"

	"github.com/avressatelier/lokan/internal/types"
	"gopkg.in/yaml.v3"
)

// parseFile reads only the frontmatter of a raw task file into a summary.
func parseFile(raw string, filePath string) (*types.TaskSummary, error) {
	fmStr, _, ok := splitFrontmatter(raw)
	if !ok {
		return nil, errNoFrontmatter
	}
	fm, err := parseFrontmatter(fmStr)
	if err != nil {
		return nil, err
	}
	return &types.TaskSummary{TaskFrontmatter: fm, FilePath: filePath}, nil
}

// parseFullFile reads frontmatter plus the markdown body into a full task.
func parseFullFile(raw string, filePath string) (*types.Task, error) {
	fmStr, body, ok := splitFrontmatter(raw)
	if !ok {
		return nil, errNoFrontmatter
	}
	fm, err := parseFrontmatter(fmStr)
	if err != nil {
		return nil, err
	}
	return &types.Task{TaskFrontmatter: fm, Body: body, FilePath: filePath}, nil
}

// serializeTask renders a task back to markdown with YAML frontmatter.
func serializeTask(task types.Task) (string, error) {
	// copy core fields, omitting empty optionals
	fm := types.TaskFrontmatter{
		ID:       task.ID,
		Title:    task.Title,
		Type:     task.Type,
		Status:   task.Status,
		Priority: task.Priority,
		Created:  task.Created,
		Updated:  task.Updated,
	}
	if task.Parent != "" {
		fm.Parent = task.Parent
	}
	if len(task.Related) > 0 {
		fm.Related = task.Related
	}
	if len(task.Docs) > 0 {
		fm.Docs = task.Docs
	}
	if len(task.Tags) > 0 {
		fm.Tags = task.Tags
	}

	// marshal the frontmatter and wrap it around the body
	raw, err := yaml.Marshal(fm)
	if err != nil {
		return "", fmt.Errorf("serialize frontmatter: %w", err)
	}
	return fmt.Sprintf("---\n%s---\n%s", raw, task.Body), nil
}

// buildInitialBody scaffolds the default markdown body for a new task.
func buildInitialBody(title string, body string) string {
	// prepend a provided body, trimmed, then add the standard sections
	content := ""
	if body != "" {
		content = strings.TrimRight(body, " \t\n") + "\n"
	}
	return fmt.Sprintf("# %s\n\n%s\n## Notes\n\n\n## Work Log\n\n", title, content)
}

// splitFrontmatter extracts the YAML block and body from a task file.
func splitFrontmatter(raw string) (fm string, body string, ok bool) {
	// normalize CRLF, then require the leading --- fence
	raw = strings.ReplaceAll(raw, "\r\n", "\n")
	if !strings.HasPrefix(raw, "---\n") {
		return "", "", false
	}

	// cut at the closing fence
	rest := raw[len("---\n"):]
	end := strings.Index(rest, "\n---\n")
	if end < 0 {
		return "", "", false
	}
	return rest[:end], rest[end+len("\n---\n"):], true
}

// parseFrontmatter validates and converts YAML into a TaskFrontmatter.
// Optional fields are strictly typed (Issue 6); invalid files are rejected.
func parseFrontmatter(fmStr string) (types.TaskFrontmatter, error) {
	var fm types.TaskFrontmatter

	// unmarshal the yaml and confirm a mapping root
	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(fmStr), &doc); err != nil {
		return fm, fmt.Errorf("parse frontmatter yaml: %w", err)
	}
	if len(doc.Content) == 0 {
		return fm, errInvalidFrontmatter
	}
	root := doc.Content[0]
	if root.Kind != yaml.MappingNode {
		return fm, errInvalidFrontmatter
	}

	// field handlers: scalar strings, enums, and string arrays
	fields := map[string]func(*yaml.Node) error{
		"id":       func(n *yaml.Node) error { return scalarString(n, &fm.ID) },
		"title":    func(n *yaml.Node) error { return scalarString(n, &fm.Title) },
		"created":  func(n *yaml.Node) error { return scalarString(n, &fm.Created) },
		"updated":  func(n *yaml.Node) error { return scalarString(n, &fm.Updated) },
		"type":     func(n *yaml.Node) error { return enumString(n, &fm.Type, types.TaskTypes) },
		"status":   func(n *yaml.Node) error { return enumString(n, &fm.Status, types.Statuses) },
		"priority": func(n *yaml.Node) error { return enumString(n, &fm.Priority, types.Priorities) },
		"parent":   func(n *yaml.Node) error { return scalarString(n, &fm.Parent) },
		"tags":     func(n *yaml.Node) error { return stringArray(n, &fm.Tags) },
		"related":  func(n *yaml.Node) error { return stringArray(n, &fm.Related) },
		"docs":     func(n *yaml.Node) error { return stringArray(n, &fm.Docs) },
	}

	// apply handlers per key, tracking which fields were seen
	seen := make(map[string]bool, len(root.Content)/2)
	for i := 0; i+1 < len(root.Content); i += 2 {
		key := root.Content[i].Value
		seen[key] = true
		handler, ok := fields[key]
		if !ok {
			continue
		}
		if err := handler(root.Content[i+1]); err != nil {
			return fm, err
		}
	}

	// require every mandatory field to be present
	for _, required := range []string{"id", "title", "type", "status", "priority", "created", "updated"} {
		if !seen[required] {
			return fm, fmt.Errorf("%w: missing required field %q", errInvalidFrontmatter, required)
		}
	}
	return fm, nil
}

func scalarString(n *yaml.Node, out *string) error {
	if n.Kind != yaml.ScalarNode || n.Tag != "!!str" {
		return fmt.Errorf("%w: %q must be a string", errInvalidFrontmatter, n.Value)
	}
	*out = n.Value
	return nil
}

func enumString[T ~string](n *yaml.Node, out *T, allowed []T) error {
	if n.Kind != yaml.ScalarNode || n.Tag != "!!str" {
		return fmt.Errorf("%w: %q must be a string", errInvalidFrontmatter, n.Value)
	}
	for _, v := range allowed {
		if T(n.Value) == v {
			*out = T(n.Value)
			return nil
		}
	}
	return fmt.Errorf("%w: invalid value %q", errInvalidFrontmatter, n.Value)
}

func stringArray(n *yaml.Node, out *[]string) error {
	if n.Kind != yaml.SequenceNode {
		return fmt.Errorf("%w: field must be an array", errInvalidFrontmatter)
	}
	vals := make([]string, len(n.Content))
	for i, item := range n.Content {
		if item.Kind != yaml.ScalarNode || item.Tag != "!!str" {
			return fmt.Errorf("%w: array element must be a string", errInvalidFrontmatter)
		}
		vals[i] = item.Value
	}
	*out = vals
	return nil
}

var (
	errNoFrontmatter      = errors.New("no yaml frontmatter")
	errInvalidFrontmatter = errors.New("invalid task frontmatter")
)
