package store

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/avressatelier/lokan/internal/types"
	"gopkg.in/yaml.v3"
)

// parseFile reads only the frontmatter of a raw task file into a summary.
func parseFile(raw string, filePath string, statuses []types.StatusDef) (*types.TaskSummary, error) {
	fmStr, _, ok := splitFrontmatter(raw)
	if !ok {
		return nil, errNoFrontmatter
	}
	fm, err := parseFrontmatter(fmStr, statuses)
	if err != nil {
		return nil, err
	}
	return &types.TaskSummary{TaskFrontmatter: fm, FilePath: filePath}, nil
}

// parseFullFile reads frontmatter plus the markdown body into a full task.
func parseFullFile(raw string, filePath string, statuses []types.StatusDef) (*types.Task, error) {
	fmStr, body, ok := splitFrontmatter(raw)
	if !ok {
		return nil, errNoFrontmatter
	}
	fm, err := parseFrontmatter(fmStr, statuses)
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
	return fmt.Sprintf("# %s\n\n%s\n## Notes\n\n\n## Work Log\n", title, content)
}

// board document layout: one file holding every task block, grouped into
// Active (open statuses) and Archive (done/cancelled) sections.
const (
	boardHeader    = "# Kanlo Board"
	sectionActive  = "## Active"
	sectionArchive = "## Archive"
	markerPrefix   = "<!-- lokan:"
)

// isArchived reports whether a task belongs in the Archive section, per the
// configured lane set. Unknown statuses default to active.
func isArchived(status types.Status, statuses []types.StatusDef) bool {
	for _, s := range statuses {
		if s.ID == status {
			return s.Archived
		}
	}
	return false
}

// statusIDs extracts the configured lane ids for enum validation.
func statusIDs(statuses []types.StatusDef) []types.Status {
	ids := make([]types.Status, len(statuses))
	for i, s := range statuses {
		ids[i] = s.ID
	}
	return ids
}

// parseBoard splits a board document into task blocks. Each task block starts
// with a "<!-- lokan:<id> -->" marker; blocks that fail to parse are skipped
// with a warning. Section headers and blank lines between blocks are dropped.
func parseBoard(raw string, statuses []types.StatusDef) []types.Task {
	var tasks []types.Task
	var block []string
	blockID := ""
	blockStart := 0
	blockEnd := 0

	flush := func() {
		if blockID == "" {
			return
		}
		// trim leading/trailing blank lines and stray section headers
		for len(block) > 0 && strings.TrimSpace(block[0]) == "" {
			block = block[1:]
		}
		for len(block) > 0 {
			last := strings.TrimSpace(block[len(block)-1])
			if last == "" || last == sectionActive || last == sectionArchive || last == boardHeader {
				block = block[:len(block)-1]
				continue
			}
			break
		}
		parsed, err := parseFullFile(strings.Join(block, "\n"), blockID, statuses)
		if err != nil {
			log.Printf("Warning: skipping invalid task block: %s", blockID)
		} else {
			// record the block's extent in the board file (1-based lines)
			parsed.LineStart = blockStart
			parsed.LineEnd = blockEnd
			if parsed.LineEnd < parsed.LineStart {
				parsed.LineEnd = parsed.LineStart
			}
			tasks = append(tasks, *parsed)
		}
		block = nil
		blockID = ""
	}

	for i, line := range strings.Split(raw, "\n") {
		if id, ok := markerID(line); ok {
			flush()
			blockID = id
			blockStart = i + 1
			blockEnd = i + 1
			continue
		}
		if blockID != "" {
			block = append(block, line)
			// extend the end only over content lines, so the reported range
			// stays tight (no trailing blanks or section headers)
			if strings.TrimSpace(line) != "" && !isSectionHeader(line) {
				blockEnd = i + 1
			}
		}
	}
	flush()
	return tasks
}

// isSectionHeader reports whether a line is a board layout header, which is
// trimmed from blocks and excluded from their reported line ranges.
func isSectionHeader(line string) bool {
	line = strings.TrimSpace(line)
	return line == sectionActive || line == sectionArchive || line == boardHeader
}

// serializeBoard renders tasks back to a single board document, grouping
// active tasks under "## Active" and finished ones under "## Archive".
func serializeBoard(tasks []types.Task, statuses []types.StatusDef) (string, error) {
	var active, archived []types.Task
	for _, t := range tasks {
		if isArchived(t.Status, statuses) {
			archived = append(archived, t)
		} else {
			active = append(active, t)
		}
	}

	var b strings.Builder
	b.WriteString(boardHeader + "\n\n")
	writeSection := func(title string, section []types.Task) error {
		b.WriteString(title + "\n")
		if len(section) == 0 {
			b.WriteString("\n")
			return nil
		}
		for _, t := range section {
			raw, err := serializeTask(t)
			if err != nil {
				return err
			}
			fmt.Fprintf(&b, "\n<!-- lokan:%s -->\n%s", t.ID, raw)
		}
		return nil
	}
	if err := writeSection(sectionActive, active); err != nil {
		return "", err
	}
	if err := writeSection(sectionArchive, archived); err != nil {
		return "", err
	}
	return b.String(), nil
}

// markerID extracts the task id from a "<!-- lokan:<id> -->" marker line.
func markerID(line string) (string, bool) {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, markerPrefix) || !strings.HasSuffix(line, " -->") {
		return "", false
	}
	id := strings.TrimSuffix(strings.TrimPrefix(line, markerPrefix), " -->")
	if id == "" {
		return "", false
	}
	return id, true
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
func parseFrontmatter(fmStr string, statuses []types.StatusDef) (types.TaskFrontmatter, error) {
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
		"status":   func(n *yaml.Node) error { return enumString(n, &fm.Status, statusIDs(statuses)) },
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
