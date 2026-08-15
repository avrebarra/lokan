package store

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/avrebarra/lokan/internal/types"
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

	// marshal the frontmatter; the marker line opens the html comment and this
	// closer hides the yaml in rendered output. The yaml is written fenceless
	// (no --- delimiters) so formatters like prettier 2 treat the comment as
	// opaque — fences and column-0 list markers get re-parsed as markdown and
	// leak the markup into rendered output.
	raw, err := yaml.Marshal(fm)
	if err != nil {
		return "", fmt.Errorf("serialize frontmatter: %w", err)
	}
	return fmt.Sprintf("%s%s\n%s", raw, commentClose, task.Body), nil
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
	boardHeader      = "# Lokan Board"
	sectionActive    = "## Active"
	sectionArchive   = "## Archive"
	markerPrefix     = "<!-- lokan:"
	configMarkerID   = "config"
	configMarkerLine = markerPrefix + configMarkerID
	// comment wrapper: the marker line opens one html comment that hides the
	// engine markup in rendered output (GitHub etc.) while the raw file stays
	// parseable. YAML is written fenceless so formatters (prettier 2) keep
	// the comment opaque; older boards with --- fences or a self-closed
	// marker line are still accepted.
	commentOpen  = "<!--"
	commentClose = "-->"
	// DefaultGuideURL is where cold-start readers learn to read this file.
	DefaultGuideURL = "https://github.com/avrebarra/lokan/blob/main/docs/guides.md"
	// boardBanner opens the board: a descriptive comment so anyone finding
	// the file without lokan knowledge understands what it is, its format,
	// and where to get help. Keep it free of "-->" so it stays one comment.
	boardBanner = `This board is a lokan kanban / roadmap — created and managed by lokan,
a single-file markdown task tool (CLI + web UI).

File format: markdown with a lokan:config block and task blocks marked
lokan:<id> (YAML frontmatter + markdown body). All engine markup is
comment-wrapped, so rendered markdown shows only the human-readable part.

Prefer the lokan tool (CLI or UI) for edits — hand-editing is possible
but the engine rewrites this file atomically on every change.

Tool:        https://github.com/avrebarra/lokan
Reference:   ` + DefaultGuideURL
)

// IsArchived reports whether a task belongs in the Archive section, per the
// configured lane set. Unknown statuses default to active.
func IsArchived(status types.Status, statuses []types.StatusDef) bool {
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
		parsed, err := parseFullFile(strings.Join(block, "\n")+"\n", blockID, statuses)
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
			if id == configMarkerID {
				continue
			}
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

// serializeBoard renders tasks back to a single board document: the config
// block first, then active tasks under "## Active" and finished ones under
// "## Archive".
func serializeBoard(tasks []types.Task, cfg types.LokanConfig) (string, error) {
	var active, archived []types.Task
	for _, t := range tasks {
		if IsArchived(t.Status, cfg.Statuses) {
			archived = append(archived, t)
		} else {
			active = append(active, t)
		}
	}

	var b strings.Builder

	// descriptive banner opens the board: self-identifies lokan, the format,
	// and the reference for readers without lokan knowledge
	b.WriteString(commentOpen + "\n")
	b.WriteString(boardBanner + "\n")
	b.WriteString(commentClose + "\n\n")

	cfgRaw, err := yaml.Marshal(cfg)
	if err != nil {
		return "", fmt.Errorf("serialize config block: %w", err)
	}
	b.WriteString(configMarkerLine + "\n")
	b.Write(cfgRaw)
	b.WriteString(commentClose + "\n\n" + boardHeader + "\n\n")
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
			fmt.Fprintf(&b, "\n<!-- lokan:%s\n%s", t.ID, raw)
		}
		return nil
	}
	if err := writeSection(sectionActive, active); err != nil {
		return "", err
	}
	// separate the archive section with a blank line so the section header
	// never glues onto the last active task's body
	if !strings.HasSuffix(b.String(), "\n\n") {
		b.WriteString("\n")
	}
	if err := writeSection(sectionArchive, archived); err != nil {
		return "", err
	}
	return b.String(), nil
}

// markerID extracts the task id from a "<!-- lokan:<id> -->" marker line or
// the merged "<!-- lokan:<id>" block opener.
func markerID(line string) (string, bool) {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, markerPrefix) {
		return "", false
	}
	id := strings.TrimSuffix(strings.TrimPrefix(line, markerPrefix), " -->")
	if id == "" {
		return "", false
	}
	return id, true
}

// isConfigMarker reports whether a line opens the board's config block
// (merged "<!-- lokan:config" or legacy self-closed "<!-- lokan:config -->").
func isConfigMarker(line string) bool {
	line = strings.TrimSpace(line)
	return line == configMarkerLine || line == configMarkerLine+" -->"
}

// splitFrontmatter extracts the YAML block and body from a task file.
func splitFrontmatter(raw string) (fm string, body string, ok bool) {
	// normalize CRLF, then unwrap the v1 comment opener (older boards carry
	// bare --- fences and are accepted unchanged)
	raw = strings.ReplaceAll(raw, "\r\n", "\n")
	if strings.HasPrefix(raw, commentOpen+"\n") {
		raw = strings.TrimPrefix(raw, commentOpen+"\n")
	}

	// legacy: fenced frontmatter, cut at the closing fence
	if strings.HasPrefix(raw, "---\n") {
		rest := raw[len("---\n"):]
		end := strings.Index(rest, "\n---\n")
		if end < 0 {
			return "", "", false
		}
		body = rest[end+len("\n---\n"):]
		// drop the comment close that follows the fence in wrapped blocks
		if strings.HasPrefix(body, commentClose+"\n") {
			body = body[len(commentClose)+1:]
		}
		return rest[:end], body, true
	}

	// current: fenceless yaml inside the comment, cut at the comment close
	end := strings.Index(raw, "\n"+commentClose+"\n")
	if end < 0 {
		return "", "", false
	}
	return raw[:end], raw[end+len("\n"+commentClose+"\n"):], true
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

// parseConfigBlock extracts the lokan config from the board's config block
// (the first "<!-- lokan:config -->" block, ending at the board header). A
// missing or unparseable block yields the defaults.
func parseConfigBlock(raw string) types.LokanConfig {
	var cfg types.LokanConfig
	lines := strings.Split(raw, "\n")
	for i, line := range lines {
		if !isConfigMarker(line) {
			continue
		}
		var yamlLines []string
		for _, l := range lines[i+1:] {
			if strings.HasPrefix(l, boardHeader) {
				break
			}
			yamlLines = append(yamlLines, l)
		}
		// strip the comment wrapper (older boards without one are unchanged)
		for len(yamlLines) > 0 && strings.TrimSpace(yamlLines[len(yamlLines)-1]) == "" {
			yamlLines = yamlLines[:len(yamlLines)-1]
		}
		if len(yamlLines) > 0 && strings.TrimSpace(yamlLines[len(yamlLines)-1]) == commentClose {
			yamlLines = yamlLines[:len(yamlLines)-1]
		}
		if len(yamlLines) > 0 && strings.TrimSpace(yamlLines[0]) == commentOpen {
			yamlLines = yamlLines[1:]
		}
		_ = yaml.Unmarshal([]byte(strings.Join(yamlLines, "\n")), &cfg)
		break
	}
	if len(cfg.Statuses) == 0 {
		cfg.Statuses = types.DefaultStatusDefs()
	}
	return cfg
}

// InitialBoard renders a fresh board document: config block plus empty
// Active/Archive sections. Used by lokan init to scaffold a new project.
func InitialBoard(cfg types.LokanConfig) (string, error) {
	return serializeBoard(nil, cfg)
}
