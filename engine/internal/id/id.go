package id

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/avressatelier/lokan/internal/types"
)

const (
	configDir      = ".lokan"
	configFile     = "config.json"
	lockTimeout    = 5 * time.Second
	lockRetryDelay = 10 * time.Millisecond
)

// ErrNotInProject is returned when no .lokan config exists for the root.
var ErrNotInProject = errors.New("not a lokan project; run `lokan init` first")

// ConfigPath returns the absolute path to .lokan/config.json under root.
func ConfigPath(root string) string {
	return filepath.Join(root, configDir, configFile)
}

// ReadConfig loads the project config, or ErrNotInProject if absent.
func ReadConfig(root string) (types.LokanConfig, error) {
	var cfg types.LokanConfig
	raw, err := os.ReadFile(ConfigPath(root))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, ErrNotInProject
		}
		return cfg, err
	}
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return cfg, fmt.Errorf("parse %s: %w", ConfigPath(root), err)
	}
	return cfg, nil
}

// WriteConfig atomically persists the config via temp-file-then-rename.
func WriteConfig(root string, cfg types.LokanConfig) error {
	// ensure the config dir exists
	if err := os.MkdirAll(filepath.Dir(ConfigPath(root)), 0o755); err != nil {
		return err
	}

	// marshal the config with a trailing newline
	raw, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	raw = append(raw, '\n')

	// write to a temp file, then atomically rename into place
	tmp, err := os.CreateTemp(filepath.Dir(ConfigPath(root)), ".config-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer func() {
		tmp.Close()
		os.Remove(tmpPath)
	}()
	if _, err := tmp.Write(raw); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpPath, ConfigPath(root))
}

// NextCounter increments the config counter and returns the new value.
// The read-modify-write is guarded by an exclusive lock file so concurrent
// callers cannot hand out duplicate counters (Issue 1).
func NextCounter(root string) (int, error) {
	// acquire the exclusive lock file, retrying until timeout
	lockPath := ConfigPath(root) + ".lock"
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
			return 0, err
		}
		if time.Now().After(deadline) {
			return 0, fmt.Errorf("timed out acquiring lock %s", lockPath)
		}
		time.Sleep(lockRetryDelay)
	}

	// read-modify-write the counter under the held lock
	cfg, err := ReadConfig(root)
	if err != nil {
		return 0, err
	}
	cfg.Counter++
	if err := WriteConfig(root, cfg); err != nil {
		return 0, err
	}
	return cfg.Counter, nil
}

// GenerateSlug lowercases a title, hyphenates non-alphanumeric runs, trims
// edge hyphens, and truncates at a word boundary below 50 characters.
func GenerateSlug(title string) string {
	// normalize to a lowercase hyphenated slug
	slug := strings.ToLower(title)
	slug = nonAlnum.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	// truncate at a word boundary past 50 chars
	if len(slug) > 50 {
		slug = trailingWord.ReplaceAllString(slug[:50], "")
	}
	return slug
}

// GenerateID builds a task id from type and counter, e.g. "task-2".
func GenerateID(taskType types.TaskType, counter int) string {
	return fmt.Sprintf("%s-%d", taskType, counter)
}

// GenerateFilename builds the task file name, e.g. "task-2-setup-db.md".
func GenerateFilename(taskType types.TaskType, counter int, title string) string {
	return fmt.Sprintf("%s-%d-%s.md", taskType, counter, GenerateSlug(title))
}

var (
	nonAlnum     = regexp.MustCompile(`[^a-z0-9]+`)
	trailingWord = regexp.MustCompile(`-[^-]*$`)
)
