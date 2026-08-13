package main

import (
	"fmt"
	"strings"

	"github.com/avressatelier/lokan/internal/types"
)

// renderTable renders tasks as a compact aligned table.
func renderTable(tasks []types.TaskSummary) string {
	if len(tasks) == 0 {
		return "No tasks found."
	}

	// build rows and compute per-column widths
	headers := []string{"ID", "TYPE", "STATUS", "PRIORITY", "TITLE"}
	rows := make([][]string, len(tasks))
	for i, t := range tasks {
		rows[i] = []string{t.ID, string(t.Type), string(t.Status), string(t.Priority), t.Title}
	}
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// local helpers for padding and line assembly
	pad := func(s string, w int) string {
		if len(s) < w {
			return s + strings.Repeat(" ", w-len(s))
		}
		return s
	}
	line := func(cells []string) string {
		parts := make([]string, len(cells))
		for i, c := range cells {
			parts[i] = pad(c, widths[i])
		}
		return strings.Join(parts, "  ")
	}
	sep := make([]string, len(widths))
	for i, w := range widths {
		sep[i] = strings.Repeat("─", w)
	}

	// render header, separator, and rows
	var b strings.Builder
	b.WriteString(line(headers))
	b.WriteString("\n")
	b.WriteString(strings.Join(sep, "──"))
	b.WriteString("\n")
	for _, row := range rows {
		b.WriteString(line(row))
		b.WriteString("\n")
	}
	return strings.TrimRight(b.String(), "\n")
}

// renderTaskDetail renders a full task: frontmatter fields plus body.
func renderTaskDetail(task types.Task) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s · %s\n", task.ID, task.Type)
	fmt.Fprintf(&b, "Title:    %s\n", task.Title)
	fmt.Fprintf(&b, "Status:   %s\n", task.Status)
	fmt.Fprintf(&b, "Priority: %s\n", task.Priority)
	if task.Parent != "" {
		fmt.Fprintf(&b, "Parent:   %s\n", task.Parent)
	}
	if len(task.Related) > 0 {
		fmt.Fprintf(&b, "Related:  %s\n", strings.Join(task.Related, ", "))
	}
	if len(task.Docs) > 0 {
		fmt.Fprintf(&b, "Docs:     %s\n", strings.Join(task.Docs, ", "))
	}
	if len(task.Tags) > 0 {
		fmt.Fprintf(&b, "Tags:     %s\n", strings.Join(task.Tags, ", "))
	}
	fmt.Fprintf(&b, "Created:  %s  Updated: %s\n", task.Created, task.Updated)
	b.WriteString(strings.Repeat("─", 50) + "\n")
	b.WriteString(task.Body)
	return b.String()
}

// rowLine renders a single task as one table row.
func rowLine(t types.TaskSummary) string {
	return fmt.Sprintf("%s  %s  %s  %s  %s", t.ID, t.Type, t.Status, t.Priority, t.Title)
}

// renderMarkdownBoard renders tasks grouped by status as compact markdown,
// one line per task — the agent-friendly list mode. Groups follow the
// configured lane order.
func renderMarkdownBoard(tasks []types.TaskSummary, statuses []types.StatusDef) string {
	if len(tasks) == 0 {
		return "No tasks found."
	}

	// bucket tasks by status
	byStatus := make(map[types.Status][]types.TaskSummary)
	for _, t := range tasks {
		byStatus[t.Status] = append(byStatus[t.Status], t)
	}

	// count active vs archived per the configured lane set
	active, archived := 0, 0
	for _, t := range tasks {
		if isArchivedStatus(t.Status, statuses) {
			archived++
		} else {
			active++
		}
	}

	// render status groups in configured lane order
	var b strings.Builder
	fmt.Fprintf(&b, "# Board — %d active, %d archived\n", active, archived)
	for _, status := range statuses {
		group, ok := byStatus[status.ID]
		if !ok {
			continue
		}
		b.WriteString("\n## " + string(status.ID) + "\n")
		for _, t := range group {
			fmt.Fprintf(&b, "- %s [%s] %s\n", t.ID, t.Priority, t.Title)
		}
	}
	return strings.TrimRight(b.String(), "\n")
}
