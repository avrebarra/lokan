package types

// Status is the lifecycle state of a task.
type Status string

const (
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in-progress"
	StatusBacklog    Status = "backlog"
	StatusDone       Status = "done"
	StatusCancelled  Status = "cancelled"
)

// Allowed enum values, in contract order. Every task is a plain task — the
// type/priority/parent dimensions were removed; legacy fields in existing
// board files are tolerated on read and never written back.
var Statuses = []Status{StatusTodo, StatusInProgress, StatusBacklog, StatusDone, StatusCancelled}

// Contains reports whether v is in list. Shared by the CLI and the API.
func Contains[T comparable](list []T, v T) bool {
	for _, item := range list {
		if item == v {
			return true
		}
	}
	return false
}

// TaskFrontmatter is the YAML frontmatter of a task file. Field order mirrors
// the reference serializer output.
type TaskFrontmatter struct {
	ID      string   `json:"id" yaml:"id"`
	Title   string   `json:"title" yaml:"title"`
	Status  Status   `json:"status" yaml:"status"`
	Created string   `json:"created" yaml:"created"`
	Updated string   `json:"updated" yaml:"updated"`
	Related []string `json:"related,omitempty" yaml:"related,omitempty"`
	Docs    []string `json:"docs,omitempty" yaml:"docs,omitempty"`
	Tags    []string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// Task is a full task: frontmatter plus raw markdown body and its on-disk path.
type Task struct {
	TaskFrontmatter
	Body      string `json:"body"`
	FilePath  string `json:"filePath"`
	LineStart int    `json:"lineStart"`
	LineEnd   int    `json:"lineEnd"`
}

// TaskSummary is frontmatter plus path, without the body.
type TaskSummary struct {
	TaskFrontmatter
	FilePath  string `json:"filePath"`
	LineStart int    `json:"lineStart"`
	LineEnd   int    `json:"lineEnd"`
}

// StatusDef defines a configurable lane on the board: its id (stored in
// task frontmatter and rendered as the column header) and whether the lane
// counts as archived (drives the Active/Archive board split and clear ops).
type StatusDef struct {
	ID       Status `json:"id" yaml:"id"`
	Archived bool   `json:"archived" yaml:"archived,omitempty"`
}

// DefaultStatusDefs returns the built-in lane set, used when a project has
// no configured lanes. Order is the contract order.
func DefaultStatusDefs() []StatusDef {
	return []StatusDef{
		{ID: StatusBacklog},
		{ID: StatusTodo},
		{ID: StatusInProgress},
		{ID: StatusDone, Archived: true},
		{ID: StatusCancelled, Archived: true},
	}
}

// LokanConfig is the config block stored in the board file.
type LokanConfig struct {
	Title    string      `json:"title,omitempty" yaml:"title,omitempty"`
	Counter  int         `json:"counter" yaml:"counter"`
	Version  string      `json:"version" yaml:"version"`
	Statuses []StatusDef `json:"statuses,omitempty" yaml:"statuses,omitempty"`
}

// QueryOptions filters tasks by the given dimensions. Zero values are ignored.
type QueryOptions struct {
	Status Status
	Tags   []string
}
