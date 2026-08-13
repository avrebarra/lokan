package store

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/avressatelier/lokan/internal/id"
	"github.com/avressatelier/lokan/internal/types"
)

const (
	tasksDirName   = ".lokan"
	boardFileName  = "board.md"
	lockTimeout    = 5 * time.Second
	lockRetryDelay = 10 * time.Millisecond
)

// ErrNotFound is returned when no task block matches a requested id.
var ErrNotFound = errors.New("task not found")

// TasksDir returns the virtual tasks directory under root. Tasks no longer
// live as separate files; the path only backs the filePath values reported
// through the API.
func TasksDir(root string) string {
	return filepath.Join(root, tasksDirName, "tasks")
}

// BoardPath returns the absolute path of the single board file under root.
func BoardPath(root string) string {
	return filepath.Join(root, tasksDirName, boardFileName)
}

// VirtualPath returns the virtual per-task path reported through the API,
// e.g. "<root>/.lokan/tasks/task-1.md". It resolves to a block in board.md.
func VirtualPath(root, id string) string {
	return filepath.Join(TasksDir(root), id+".md")
}

// boardPath resolves the board file path from a virtual task path.
func boardPath(virtualPath string) string {
	return filepath.Join(filepath.Dir(filepath.Dir(virtualPath)), boardFileName)
}

// rootFromVirtual resolves the project root from a virtual task path.
func rootFromVirtual(virtualPath string) string {
	return filepath.Dir(filepath.Dir(filepath.Dir(virtualPath)))
}

// statusDefs resolves the effective lane definitions for a project, falling
// back to the built-in defaults when the config is unreadable.
func statusDefs(root string) []types.StatusDef {
	cfg, err := id.ReadConfig(root)
	if err != nil {
		return types.DefaultStatusDefs()
	}
	return cfg.Statuses
}

// LoadAllSummaries reads every task block in the board file, skipping blocks
// that fail to parse (with a warning).
func LoadAllSummaries(root string) ([]types.TaskSummary, error) {
	tasks, err := readBoard(root)
	if err != nil {
		return nil, err
	}
	summaries := make([]types.TaskSummary, len(tasks))
	for i, t := range tasks {
		summaries[i] = types.TaskSummary{
			TaskFrontmatter: t.TaskFrontmatter,
			FilePath:        VirtualPath(root, t.ID),
			LineStart:       t.LineStart,
			LineEnd:         t.LineEnd,
		}
	}
	return summaries, nil
}

// FindByID scans the board for the task block matching id.
func FindByID(root string, id string) (types.TaskSummary, error) {
	var summary types.TaskSummary
	tasks, err := readBoard(root)
	if err != nil {
		return summary, err
	}
	for _, t := range tasks {
		if t.ID == id {
			return types.TaskSummary{
				TaskFrontmatter: t.TaskFrontmatter,
				FilePath:        VirtualPath(root, t.ID),
				LineStart:       t.LineStart,
				LineEnd:         t.LineEnd,
			}, nil
		}
	}
	return summary, fmt.Errorf("%w: %s", ErrNotFound, id)
}

// LoadTask reads and fully parses the task block addressed by a virtual path.
func LoadTask(virtualPath string) (types.Task, error) {
	var task types.Task
	id := strings.TrimSuffix(filepath.Base(virtualPath), ".md")
	tasks, err := readBoard(rootFromVirtual(virtualPath))
	if err != nil {
		return task, err
	}
	for _, t := range tasks {
		if t.ID == id {
			t.FilePath = virtualPath
			return t, nil
		}
	}
	return task, fmt.Errorf("%w: %s", ErrNotFound, id)
}

// WriteTask persists a task by replacing its block in the board, bumping the
// updated field to today (UTC). The mutation is lock-guarded and committed
// atomically via temp-file-then-rename.
func WriteTask(task types.Task) error {
	task.Updated = time.Now().UTC().Format("2006-01-02")
	return withBoardLock(boardPath(task.FilePath), func() error {
		root := rootFromVirtual(task.FilePath)
		tasks, err := readBoard(root)
		if err != nil {
			return err
		}
		replaced := false
		for i, t := range tasks {
			if t.ID == task.ID {
				tasks[i] = task
				replaced = true
				break
			}
		}
		if !replaced {
			tasks = append(tasks, task)
		}
		return writeBoard(root, tasks)
	})
}

// MoveTask relocates a task to another lane (status) and position: directly
// before beforeID, or at the end of the lane when beforeID is empty. The
// whole board is rewritten under one lock, so the move is atomic.
func MoveTask(root, id string, status types.Status, beforeID string) (types.Task, error) {
	var moved types.Task
	err := withBoardLock(BoardPath(root), func() error {
		tasks, err := readBoard(root)
		if err != nil {
			return err
		}
		// locate the task to move and pull it out of the board order
		from := -1
		for i, t := range tasks {
			if t.ID == id {
				from = i
				moved = t
				break
			}
		}
		if from < 0 {
			return fmt.Errorf("%w: %s", ErrNotFound, id)
		}
		tasks = append(tasks[:from], tasks[from+1:]...)

		// apply the move: reinsert at the anchor, then write once
		moved.Status = status
		moved.Updated = time.Now().UTC().Format("2006-01-02")
		to := insertIndexFor(tasks, status, beforeID, statusDefs(root))
		tasks = append(tasks[:to], append([]types.Task{moved}, tasks[to:]...)...)
		return writeBoard(root, tasks)
	})
	if err != nil {
		return moved, err
	}
	return moved, nil
}

// insertIndexFor computes the board position for a task landing in the given
// lane: before the anchor task, after the lane's last task, or at the end of
// the active/archive section when the lane is empty.
func insertIndexFor(tasks []types.Task, status types.Status, beforeID string, statuses []types.StatusDef) int {
	if beforeID != "" {
		for i, t := range tasks {
			if t.ID == beforeID {
				return i
			}
		}
		return len(tasks)
	}
	// append after the lane's last task, keeping lane order stable
	for i := len(tasks) - 1; i >= 0; i-- {
		if tasks[i].Status == status {
			return i + 1
		}
	}
	// empty lane: append at the end of the matching section
	if isArchived(status, statuses) {
		return len(tasks)
	}
	for i, t := range tasks {
		if isArchived(t.Status, statuses) {
			return i
		}
	}
	return len(tasks)
}

// CreateTask appends a new task block to the board and returns it.
func CreateTask(root string, fm types.TaskFrontmatter, body string) (types.Task, error) {
	var task types.Task
	if err := os.MkdirAll(filepath.Dir(BoardPath(root)), 0o755); err != nil {
		return task, err
	}
	task = types.Task{
		TaskFrontmatter: fm,
		Body:            buildInitialBody(fm.Title, body),
		FilePath:        VirtualPath(root, fm.ID),
	}
	err := withBoardLock(BoardPath(root), func() error {
		tasks, err := readBoard(root)
		if err != nil {
			return err
		}
		tasks = append(tasks, task)
		return writeBoard(root, tasks)
	})
	if err != nil {
		return task, err
	}
	return task, nil
}

// MoveLane rewrites the stored status of every task in the from lane to the
// to lane — one atomic board rewrite. Used for lane renames and for moving
// a removed lane's tasks into another lane. Returns how many tasks moved.
func MoveLane(root string, from, to types.Status) (int, error) {
	moved := 0
	err := withBoardLock(BoardPath(root), func() error {
		tasks, err := readBoard(root)
		if err != nil {
			return err
		}
		for i := range tasks {
			if tasks[i].Status == from {
				tasks[i].Status = to
				moved++
			}
		}
		return writeBoard(root, tasks)
	})
	return moved, err
}

// ClearArchived deletes every task whose lane is marked archived, returning
// how many were removed.
func ClearArchived(root string) (int, error) {
	return clearTasks(root, func(t types.Task) bool {
		return isArchived(t.Status, statusDefs(root))
	})
}

// ClearAll deletes every task on the board, returning how many were removed.
func ClearAll(root string) (int, error) {
	return clearTasks(root, func(types.Task) bool { return true })
}

// clearTasks drops the tasks matching drop and rewrites the board atomically.
func clearTasks(root string, drop func(types.Task) bool) (int, error) {
	deleted := 0
	err := withBoardLock(BoardPath(root), func() error {
		tasks, err := readBoard(root)
		if err != nil {
			return err
		}
		kept := tasks[:0]
		for _, t := range tasks {
			if drop(t) {
				deleted++
				continue
			}
			kept = append(kept, t)
		}
		return writeBoard(root, kept)
	})
	return deleted, err
}

// readBoard loads all task blocks from the board file. A missing board file
// is treated as an empty board.
func readBoard(root string) ([]types.Task, error) {
	raw, err := os.ReadFile(BoardPath(root))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	return parseBoard(string(raw), statusDefs(root)), nil
}

// writeBoard atomically persists the full board document.
func writeBoard(root string, tasks []types.Task) error {
	raw, err := serializeBoard(tasks, statusDefs(root))
	if err != nil {
		return err
	}
	board := BoardPath(root)
	tmp, err := os.CreateTemp(filepath.Dir(board), ".board-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer func() {
		tmp.Close()
		os.Remove(tmpPath)
	}()
	if _, err := tmp.WriteString(raw); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpPath, board)
}

// withBoardLock serializes board mutations with an exclusive lock file so
// concurrent writers cannot clobber each other's full-document rewrites.
func withBoardLock(board string, fn func() error) error {
	lockPath := board + ".lock"
	deadline := time.Now().Add(lockTimeout)
	for {
		lock, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
		if err == nil {
			defer func() {
				lock.Close()
				os.Remove(lockPath)
			}()
			break
		}
		if !errors.Is(err, os.ErrExist) {
			return err
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("timed out acquiring lock %s", lockPath)
		}
		time.Sleep(lockRetryDelay)
	}
	return fn()
}
