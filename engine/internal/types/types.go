package types

// TaskType is the kind of task, mirroring the frozen API contract.
type TaskType string

const (
	TypeEpic    TaskType = "epic"
	TypeTask    TaskType = "task"
	TypeSubtask TaskType = "subtask"
	TypeBug     TaskType = "bug"
)

// Status is the lifecycle state of a task.
type Status string

const (
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in-progress"
	StatusBacklog    Status = "backlog"
	StatusDone       Status = "done"
	StatusCancelled  Status = "cancelled"
)

// Priority is the urgency of a task.
type Priority string

const (
	PriorityCritical Priority = "critical"
	PriorityHigh     Priority = "high"
	PriorityMedium   Priority = "medium"
	PriorityLow      Priority = "low"
)

// Allowed enum values, in contract order.
var (
	TaskTypes  = []TaskType{TypeEpic, TypeTask, TypeSubtask, TypeBug}
	Statuses   = []Status{StatusTodo, StatusInProgress, StatusBacklog, StatusDone, StatusCancelled}
	Priorities = []Priority{PriorityCritical, PriorityHigh, PriorityMedium, PriorityLow}
)

// AllowedParents maps each task type to the task types allowed as its parent.
var AllowedParents = map[TaskType][]TaskType{
	TypeEpic:    {},
	TypeTask:    {TypeEpic},
	TypeSubtask: {TypeTask, TypeBug},
	TypeBug:     {TypeEpic, TypeTask},
}

// TaskFrontmatter is the YAML frontmatter of a task file. Field order mirrors
// the reference serializer output.
type TaskFrontmatter struct {
	ID       string   `json:"id" yaml:"id"`
	Title    string   `json:"title" yaml:"title"`
	Type     TaskType `json:"type" yaml:"type"`
	Status   Status   `json:"status" yaml:"status"`
	Priority Priority `json:"priority" yaml:"priority"`
	Created  string   `json:"created" yaml:"created"`
	Updated  string   `json:"updated" yaml:"updated"`
	Parent   string   `json:"parent,omitempty" yaml:"parent,omitempty"`
	Related  []string `json:"related,omitempty" yaml:"related,omitempty"`
	Docs     []string `json:"docs,omitempty" yaml:"docs,omitempty"`
	Tags     []string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// Task is a full task: frontmatter plus raw markdown body and its on-disk path.
type Task struct {
	TaskFrontmatter
	Body     string `json:"body"`
	FilePath string `json:"filePath"`
}

// TaskSummary is frontmatter plus path, without the body.
type TaskSummary struct {
	TaskFrontmatter
	FilePath string `json:"filePath"`
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
	Counter  int         `json:"counter" yaml:"counter"`
	Version  string      `json:"version" yaml:"version"`
	Statuses []StatusDef `json:"statuses,omitempty" yaml:"statuses,omitempty"`
}

// QueryOptions filters tasks by the given dimensions. Zero values are ignored.
type QueryOptions struct {
	Type     TaskType
	Status   Status
	Priority Priority
	Parent   string
	Tags     []string
}

// TreeNode is a task with its nested children.
type TreeNode struct {
	Task     TaskSummary
	Children []TreeNode
}
