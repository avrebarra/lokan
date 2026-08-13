package store

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/avressatelier/lokan/internal/types"
)

const (
	lockTimeout    = 5 * time.Second
	lockRetryDelay = 10 * time.Millisecond
)

// ErrNotFound is returned when no task block matches a requested id.
var ErrNotFound = errors.New("task not found")

// IsBoard reports whether path exists and its first non-blank line is the
// lokan config marker — the only thing that makes a file a board.
func IsBoard(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if line := strings.TrimSpace(sc.Text()); line != "" {
			return line == configMarkerLine
		}
	}
	return false
}

// VirtualPath returns the stable virtual task path reported through the API,
// e.g. "<board>#1". It addresses a block inside the board file and is not a
// real file on disk.
func VirtualPath(board, id string) string {
	return board + "#" + id
}

// boardFromVirtual extracts the board file path from a virtual task path.
func boardFromVirtual(virtualPath string) string {
	board, _, _ := strings.Cut(virtualPath, "#")
	return board
}

// ReadConfig loads the board's config block; a board with no config block
// yields the defaults.
func ReadConfig(board string) (types.LokanConfig, error) {
	raw, err := os.ReadFile(board)
	if err != nil {
		return types.LokanConfig{}, err
	}
	return parseConfigBlock(string(raw)), nil
}

// defaultConfig is the config a fresh board starts with.
func defaultConfig() types.LokanConfig {
	return types.LokanConfig{Counter: 0, Version: "1", Statuses: types.DefaultStatusDefs()}
}

// WriteConfig rewrites the board's config block, preserving all tasks.
func WriteConfig(board string, cfg types.LokanConfig) error {
	return withBoardLock(board, func() error {
		tasks, _, err := readBoard(board)
		if err != nil {
			return err
		}
		return writeBoard(board, tasks, cfg)
	})
}

// NextCounter increments the counter in the board's config block and returns
// the new value, guarded by the board lock so concurrent callers never hand
// out duplicate ids.
func NextCounter(board string) (int, error) {
	var next int
	err := withBoardLock(board, func() error {
		tasks, cfg, err := readBoard(board)
		if err != nil {
			return err
		}
		cfg.Counter++
		next = cfg.Counter
		return writeBoard(board, tasks, cfg)
	})
	if err != nil {
		return 0, err
	}
	return next, nil
}

// statusDefs resolves the effective lane definitions for a board, falling
// back to the built-in defaults when the board carries none.
func statusDefs(board string) []types.StatusDef {
	cfg, err := ReadConfig(board)
	if err != nil {
		return types.DefaultStatusDefs()
	}
	return cfg.Statuses
}

// LoadAllSummaries reads every task block in the board file, skipping blocks
// that fail to parse (with a warning).
func LoadAllSummaries(board string) ([]types.TaskSummary, error) {
	tasks, _, err := readBoard(board)
	if err != nil {
		return nil, err
	}
	summaries := make([]types.TaskSummary, len(tasks))
	for i, t := range tasks {
		summaries[i] = types.TaskSummary{TaskFrontmatter: t.TaskFrontmatter, FilePath: VirtualPath(board, t.ID)}
	}
	return summaries, nil
}

// FindByID scans the board for the task block matching id.
func FindByID(board string, id string) (types.TaskSummary, error) {
	var summary types.TaskSummary
	tasks, _, err := readBoard(board)
	if err != nil {
		return summary, err
	}
	for _, t := range tasks {
		if t.ID == id {
			return types.TaskSummary{TaskFrontmatter: t.TaskFrontmatter, FilePath: VirtualPath(board, t.ID)}, nil
		}
	}
	return summary, fmt.Errorf("%w: %s", ErrNotFound, id)
}

// LoadTask reads and fully parses the task block addressed by a virtual path.
func LoadTask(virtualPath string) (types.Task, error) {
	var task types.Task
	board, id, _ := strings.Cut(virtualPath, "#")
	tasks, _, err := readBoard(board)
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
	board := boardFromVirtual(task.FilePath)
	return withBoardLock(board, func() error {
		tasks, cfg, err := readBoard(board)
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
		return writeBoard(board, tasks, cfg)
	})
}

// MoveTask relocates a task to another lane (status) and position: directly
// before beforeID, or at the end of the lane when beforeID is empty. The
// whole board is rewritten under one lock, so the move is atomic.
func MoveTask(board string, id string, status types.Status, beforeID string) (types.Task, error) {
	var moved types.Task
	err := withBoardLock(board, func() error {
		tasks, cfg, err := readBoard(board)
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
		to := insertIndexFor(tasks, status, beforeID, cfg.Statuses)
		tasks = append(tasks[:to], append([]types.Task{moved}, tasks[to:]...)...)
		return writeBoard(board, tasks, cfg)
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
func CreateTask(board string, fm types.TaskFrontmatter, body string) (types.Task, error) {
	task := types.Task{
		TaskFrontmatter: fm,
		Body:            buildInitialBody(fm.Title, body),
		FilePath:        VirtualPath(board, fm.ID),
	}
	err := withBoardLock(board, func() error {
		tasks, cfg, err := readBoard(board)
		if err != nil {
			return err
		}
		tasks = append(tasks, task)
		return writeBoard(board, tasks, cfg)
	})
	if err != nil {
		return task, err
	}
	return task, nil
}

// UpdateLanes atomically applies lane renames and persists the new lane set
// in one board rewrite, so task statuses and the config block stay
// consistent (a board can never hold a task whose lane is not configured).
func UpdateLanes(board string, renames [][2]types.Status, statuses []types.StatusDef) (int, error) {
	moved := 0
	err := withBoardLock(board, func() error {
		tasks, cfg, err := readBoard(board)
		if err != nil {
			return err
		}
		for _, pair := range renames {
			for i := range tasks {
				if tasks[i].Status == pair[0] {
					tasks[i].Status = pair[1]
					moved++
				}
			}
		}
		cfg.Statuses = statuses
		return writeBoard(board, tasks, cfg)
	})
	return moved, err
}

// MoveLane rewrites the stored status of every task in the from lane to the
// to lane — one atomic board rewrite. Returns how many tasks moved.
func MoveLane(board string, from, to types.Status) (int, error) {
	moved := 0
	err := withBoardLock(board, func() error {
		tasks, cfg, err := readBoard(board)
		if err != nil {
			return err
		}
		for i := range tasks {
			if tasks[i].Status == from {
				tasks[i].Status = to
				moved++
			}
		}
		return writeBoard(board, tasks, cfg)
	})
	return moved, err
}

// ClearArchived deletes every task whose lane is marked archived, returning
// how many were removed.
func ClearArchived(board string) (int, error) {
	return clearTasks(board, func(t types.Task) bool {
		return isArchived(t.Status, statusDefs(board))
	})
}

// ClearAll deletes every task on the board, returning how many were removed.
func ClearAll(board string) (int, error) {
	return clearTasks(board, func(types.Task) bool { return true })
}

// clearTasks drops the tasks matching drop and rewrites the board atomically.
func clearTasks(board string, drop func(types.Task) bool) (int, error) {
	deleted := 0
	err := withBoardLock(board, func() error {
		tasks, cfg, err := readBoard(board)
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
		return writeBoard(board, kept, cfg)
	})
	return deleted, err
}

// readBoard loads the board file, returning its task blocks and config block.
func readBoard(board string) ([]types.Task, types.LokanConfig, error) {
	raw, err := os.ReadFile(board)
	if err != nil {
		return nil, types.LokanConfig{}, err
	}
	cfg := parseConfigBlock(string(raw))
	return parseBoard(string(raw), cfg.Statuses), cfg, nil
}

// writeBoard atomically persists the full board document: config block first,
// then tasks grouped into Active/Archive sections.
func writeBoard(board string, tasks []types.Task, cfg types.LokanConfig) error {
	raw, err := serializeBoard(tasks, cfg)
	if err != nil {
		return err
	}
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
