package store

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/avressatelier/lokan/internal/types"
)

const tasksDirName = ".lokan"

// ErrNotFound is returned when no task file matches a requested id.
var ErrNotFound = errors.New("task not found")

// TasksDir returns the absolute tasks directory under root.
func TasksDir(root string) string {
	return filepath.Join(root, tasksDirName, "tasks")
}

// LoadAllSummaries reads every task file in the tasks dir, skipping files
// that fail to parse (with a warning).
func LoadAllSummaries(root string) ([]types.TaskSummary, error) {
	dir := TasksDir(root)

	// read the tasks dir, tolerating a missing dir as empty
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	// parse each .md file, skipping invalid ones
	var summaries []types.TaskSummary
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		filePath := filepath.Join(dir, entry.Name())
		summary, err := loadSummaryFile(filePath)
		if err != nil {
			log.Printf("Warning: skipping invalid task file: %s", entry.Name())
			continue
		}
		summaries = append(summaries, summary)
	}
	return summaries, nil
}

// LoadTask reads and fully parses a single task file.
func LoadTask(filePath string) (types.Task, error) {
	var task types.Task
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return task, err
	}
	parsed, err := parseFullFile(string(raw), filePath)
	if err != nil {
		return task, fmt.Errorf("failed to parse task file %s: %w", filePath, err)
	}
	return *parsed, nil
}

// FindByID resolves a task summary by scanning only for the id's file prefix
// (id + "-*.md") instead of parsing the whole tasks dir (Issue 5).
func FindByID(root string, id string) (types.TaskSummary, error) {
	var summary types.TaskSummary
	matches, err := filepath.Glob(filepath.Join(TasksDir(root), id+"-*.md"))
	if err != nil {
		return summary, err
	}
	if len(matches) == 0 {
		return summary, fmt.Errorf("%w: %s", ErrNotFound, id)
	}
	return loadSummaryFile(matches[0])
}

// WriteTask persists a task, bumping the updated field to today (UTC).
func WriteTask(task types.Task) error {
	task.Updated = time.Now().UTC().Format("2006-01-02")
	raw, err := serializeTask(task)
	if err != nil {
		return err
	}
	return os.WriteFile(task.FilePath, []byte(raw), 0o644)
}

// RenameTask moves a task to a new filename in the same directory and returns
// the task with its new path. The original file is removed.
func RenameTask(task types.Task, newFilename string) (types.Task, error) {
	dir := filepath.Dir(task.FilePath)
	newPath := filepath.Join(dir, newFilename)

	// serialize the task to its new path, then remove the old file
	updated := task
	updated.FilePath = newPath
	raw, err := serializeTask(updated)
	if err != nil {
		return task, err
	}
	if err := os.WriteFile(newPath, []byte(raw), 0o644); err != nil {
		return task, err
	}
	if err := os.Remove(task.FilePath); err != nil {
		return task, err
	}
	return updated, nil
}

// CreateTask writes a new task file with the given name and returns it.
func CreateTask(root string, fm types.TaskFrontmatter, filename string, body string) (types.Task, error) {
	var task types.Task
	dir := TasksDir(root)

	// ensure the tasks dir, assemble the task, serialize it
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return task, err
	}
	task = types.Task{
		TaskFrontmatter: fm,
		Body:            buildInitialBody(fm.Title, body),
		FilePath:        filepath.Join(dir, filename),
	}
	raw, err := serializeTask(task)
	if err != nil {
		return task, err
	}

	// write the task file
	if err := os.WriteFile(task.FilePath, []byte(raw), 0o644); err != nil {
		return task, err
	}
	return task, nil
}

func loadSummaryFile(filePath string) (types.TaskSummary, error) {
	var summary types.TaskSummary
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return summary, err
	}
	parsed, err := parseFile(string(raw), filePath)
	if err != nil {
		return summary, err
	}
	return *parsed, nil
}
