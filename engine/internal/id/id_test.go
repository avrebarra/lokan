package id

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/avressatelier/lokan/internal/types"
)

func newTempProject(t *testing.T, counter int) string {
	t.Helper()
	root := t.TempDir()
	cfgDir := filepath.Join(root, configDir)
	if err := os.MkdirAll(filepath.Join(cfgDir, "tasks"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeTestConfig(t, root, types.LokanConfig{Counter: counter, Version: "1"})
	return root
}

func writeTestConfig(t *testing.T, root string, cfg types.LokanConfig) {
	t.Helper()
	raw, err := json.Marshal(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(ConfigPath(root), raw, 0o644); err != nil {
		t.Fatal(err)
	}
}

// ---------------------------------------------------------------------------
// generateSlug
// ---------------------------------------------------------------------------

func TestGenerateSlugLowercasesAndHyphenates(t *testing.T) {
	if got := GenerateSlug("Hello World"); got != "hello-world" {
		t.Fatalf("GenerateSlug(Hello World) = %q, want %q", got, "hello-world")
	}
}

func TestGenerateSlugCollapsesHyphensAndTrims(t *testing.T) {
	if got := GenerateSlug("  --Multiple---Hyphens-- "); got != "multiple-hyphens" {
		t.Fatalf("got %q", got)
	}
}

func TestGenerateSlugTruncatesAtWordBoundary(t *testing.T) {
	long := "This is a very long title that definitely exceeds fifty characters in total"
	got := GenerateSlug(long)
	if len(got) > 50 {
		t.Fatalf("slug too long: %d", len(got))
	}
	if len(got) > 0 && got[len(got)-1] == '-' {
		t.Fatalf("slug ends mid-word: %q", got)
	}
	fullSlug := "this-is-a-very-long-title-that-definitely-exceeds-fifty-characters-in-total"
	if !strings.HasPrefix(fullSlug, got) {
		t.Fatalf("slug %q is not a prefix of %q", got, fullSlug)
	}
}

func TestGenerateSlugSpecialOnly(t *testing.T) {
	if got := GenerateSlug("!!!"); got != "" {
		t.Fatalf("got %q, want empty", got)
	}
	if got := GenerateSlug(""); got != "" {
		t.Fatalf("got %q, want empty", got)
	}
	if got := GenerateSlug("123"); got != "123" {
		t.Fatalf("got %q, want 123", got)
	}
}

// ---------------------------------------------------------------------------
// generateId / generateFilename
// ---------------------------------------------------------------------------

func TestGenerateID(t *testing.T) {
	cases := map[string]struct {
		taskType types.TaskType
		counter  int
		want     string
	}{
		"task":    {types.TypeTask, 2, "task-2"},
		"epic":    {types.TypeEpic, 1, "epic-1"},
		"bug":     {types.TypeBug, 99, "bug-99"},
		"subtask": {types.TypeSubtask, 7, "subtask-7"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := GenerateID(tc.taskType, tc.counter); got != tc.want {
				t.Fatalf("GenerateID = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestGenerateFilename(t *testing.T) {
	got := GenerateFilename(types.TypeTask, 2, "Setup DB")
	if got != "task-2-setup-db.md" {
		t.Fatalf("got %q", got)
	}
}

// ---------------------------------------------------------------------------
// config read/write
// ---------------------------------------------------------------------------

func TestReadConfig(t *testing.T) {
	root := newTempProject(t, 3)
	cfg, err := ReadConfig(root)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Counter != 3 || cfg.Version != "1" {
		t.Fatalf("cfg = %+v", cfg)
	}
}

func TestReadConfigNotInProject(t *testing.T) {
	root := t.TempDir()
	if _, err := ReadConfig(root); !errors.Is(err, ErrNotInProject) {
		t.Fatalf("err = %v, want ErrNotInProject", err)
	}
}

func TestWriteConfigAtomically(t *testing.T) {
	root := newTempProject(t, 0)
	if err := WriteConfig(root, types.LokanConfig{Counter: 41, Version: "2"}); err != nil {
		t.Fatal(err)
	}
	cfg, err := ReadConfig(root)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Counter != 41 || cfg.Version != "2" {
		t.Fatalf("cfg = %+v", cfg)
	}
}

// ---------------------------------------------------------------------------
// nextCounter
// ---------------------------------------------------------------------------

func TestNextCounterIncrements(t *testing.T) {
	root := newTempProject(t, 0)
	if got := mustNextCounter(t, root); got != 1 {
		t.Fatalf("first = %d, want 1", got)
	}
	if got := mustNextCounter(t, root); got != 2 {
		t.Fatalf("second = %d, want 2", got)
	}
	if got := mustNextCounter(t, root); got != 3 {
		t.Fatalf("third = %d, want 3", got)
	}
}

func TestNextCounterConcurrentUnique(t *testing.T) {
	root := newTempProject(t, 0)
	const n = 50

	var wg sync.WaitGroup
	counts := make(chan int, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c, err := NextCounter(root)
			if err != nil {
				t.Errorf("NextCounter: %v", err)
				return
			}
			counts <- c
		}()
	}
	wg.Wait()
	close(counts)

	seen := make(map[int]bool, n)
	for c := range counts {
		if seen[c] {
			t.Fatalf("duplicate counter handed out: %d", c)
		}
		if c < 1 || c > n {
			t.Fatalf("counter out of range: %d", c)
		}
		seen[c] = true
	}
	if len(seen) != n {
		t.Fatalf("got %d unique counters, want %d", len(seen), n)
	}

	cfg, err := ReadConfig(root)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Counter != n {
		t.Fatalf("config counter = %d, want %d", cfg.Counter, n)
	}
}

func TestNextCounterRemovesLockFile(t *testing.T) {
	root := newTempProject(t, 0)
	mustNextCounter(t, root)
	lockPath := ConfigPath(root) + ".lock"
	if _, err := os.Stat(lockPath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("lock file still exists: %v", err)
	}
}

func mustNextCounter(t *testing.T, root string) int {
	t.Helper()
	c, err := NextCounter(root)
	if err != nil {
		t.Fatal(err)
	}
	return c
}
