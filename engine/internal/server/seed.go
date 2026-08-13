package server

import (
	"time"

	"github.com/avressatelier/lokan/internal/id"
	"github.com/avressatelier/lokan/internal/store"
	"github.com/avressatelier/lokan/internal/types"
)

// SeedDemoData creates the demo airline tasks and returns how many tasks were
// created. IDs come from the project counter.
func SeedDemoData(root string) (int, error) {
	created := 0

	// epic 1: new route launch
	epic1, err := seedTask(root, types.TaskFrontmatter{
		Type: types.TypeEpic, Title: "Launch SFO–NRT Route",
		Status: types.StatusInProgress, Priority: types.PriorityCritical,
		Tags: []string{"routes", "international"},
	})
	if err != nil {
		return created, err
	}
	created++

	t1, err := seedTask(root, types.TaskFrontmatter{
		Type: types.TypeTask, Title: "Submit route approval to FAA and JCAB",
		Status: types.StatusDone, Priority: types.PriorityCritical, Parent: epic1,
		Tags: []string{"compliance", "international"},
	})
	if err != nil {
		return created, err
	}
	created++

	if _, err := seedTask(root, types.TaskFrontmatter{
		Type: types.TypeSubtask, Title: "Prepare bilateral air service agreement documents",
		Status: types.StatusDone, Priority: types.PriorityHigh, Parent: t1,
	}); err != nil {
		return created, err
	}
	created++

	if _, err := seedTask(root, types.TaskFrontmatter{
		Type: types.TypeSubtask, Title: "Coordinate slot allocation at Narita Airport",
		Status: types.StatusDone, Priority: types.PriorityHigh, Parent: t1,
	}); err != nil {
		return created, err
	}
	created++

	t2, err := seedTask(root, types.TaskFrontmatter{
		Type: types.TypeTask, Title: "Assign Boeing 787-9 fleet for long-haul",
		Status: types.StatusInProgress, Priority: types.PriorityHigh, Parent: epic1,
		Tags: []string{"fleet", "operations"},
	})
	if err != nil {
		return created, err
	}
	created++

	if _, err := seedTask(root, types.TaskFrontmatter{
		Type: types.TypeSubtask, Title: "Schedule pre-departure maintenance check",
		Status: types.StatusInProgress, Priority: types.PriorityCritical, Parent: t2,
	}); err != nil {
		return created, err
	}
	created++

	if _, err := seedTask(root, types.TaskFrontmatter{
		Type: types.TypeSubtask, Title: "Configure cabin for 14-hour flight (meals, IFE)",
		Status: types.StatusTodo, Priority: types.PriorityMedium, Parent: t2,
	}); err != nil {
		return created, err
	}
	created++

	if _, err := seedTask(root, types.TaskFrontmatter{
		Type: types.TypeTask, Title: "Hire and train Japan-route cabin crew",
		Status: types.StatusTodo, Priority: types.PriorityHigh, Parent: epic1,
		Tags: []string{"crew", "training"},
	}); err != nil {
		return created, err
	}
	created++

	if _, err := seedTask(root, types.TaskFrontmatter{
		Type: types.TypeTask, Title: "Launch marketing campaign for SFO–NRT",
		Status: types.StatusBacklog, Priority: types.PriorityMedium, Parent: epic1,
		Tags: []string{"marketing"},
	}); err != nil {
		return created, err
	}
	created++

	// epic 2: passenger experience upgrade
	epic2, err := seedTask(root, types.TaskFrontmatter{
		Type: types.TypeEpic, Title: "Upgrade Passenger Experience",
		Status: types.StatusInProgress, Priority: types.PriorityHigh,
		Tags: []string{"passenger", "experience"},
	})
	if err != nil {
		return created, err
	}
	created++

	t3, err := seedTask(root, types.TaskFrontmatter{
		Type: types.TypeTask, Title: "Redesign business class seating",
		Status: types.StatusDone, Priority: types.PriorityHigh, Parent: epic2,
		Tags: []string{"cabin", "design"},
	})
	if err != nil {
		return created, err
	}
	created++

	if _, err := seedTask(root, types.TaskFrontmatter{
		Type: types.TypeSubtask, Title: "Source lie-flat seat suppliers",
		Status: types.StatusDone, Priority: types.PriorityHigh, Parent: t3,
	}); err != nil {
		return created, err
	}
	created++

	if _, err := seedTask(root, types.TaskFrontmatter{
		Type: types.TypeSubtask, Title: "Conduct passenger comfort testing",
		Status: types.StatusDone, Priority: types.PriorityMedium, Parent: t3,
	}); err != nil {
		return created, err
	}
	created++

	t4, err := seedTask(root, types.TaskFrontmatter{
		Type: types.TypeTask, Title: "Roll out in-flight Wi-Fi on all A320s",
		Status: types.StatusInProgress, Priority: types.PriorityHigh, Parent: epic2,
		Tags: []string{"connectivity", "fleet"},
	})
	if err != nil {
		return created, err
	}
	created++

	if _, err := seedTask(root, types.TaskFrontmatter{
		Type: types.TypeSubtask, Title: "Negotiate Starlink aviation contract",
		Status: types.StatusDone, Priority: types.PriorityCritical, Parent: t4,
	}); err != nil {
		return created, err
	}
	created++

	if _, err := seedTask(root, types.TaskFrontmatter{
		Type: types.TypeSubtask, Title: "Install antenna hardware on 12 aircraft",
		Status: types.StatusInProgress, Priority: types.PriorityHigh, Parent: t4,
	}); err != nil {
		return created, err
	}
	created++

	if _, err := seedTask(root, types.TaskFrontmatter{
		Type: types.TypeSubtask, Title: "Test bandwidth at cruising altitude",
		Status: types.StatusTodo, Priority: types.PriorityMedium, Parent: t4,
	}); err != nil {
		return created, err
	}
	created++

	if _, err := seedTask(root, types.TaskFrontmatter{
		Type: types.TypeTask, Title: "Introduce premium meal service with local chefs",
		Status: types.StatusTodo, Priority: types.PriorityMedium, Parent: epic2,
		Tags: []string{"catering", "passenger"},
	}); err != nil {
		return created, err
	}
	created++

	if _, err := seedTask(root, types.TaskFrontmatter{
		Type: types.TypeTask, Title: "Refresh loyalty program tier benefits",
		Status: types.StatusBacklog, Priority: types.PriorityLow, Parent: epic2,
		Tags: []string{"loyalty", "passenger"},
	}); err != nil {
		return created, err
	}
	created++

	// epic 3: fleet maintenance overhaul
	epic3, err := seedTask(root, types.TaskFrontmatter{
		Type: types.TypeEpic, Title: "Fleet Maintenance Overhaul Q3",
		Status: types.StatusTodo, Priority: types.PriorityCritical,
		Tags: []string{"maintenance", "safety"},
	})
	if err != nil {
		return created, err
	}
	created++

	if _, err := seedTask(root, types.TaskFrontmatter{
		Type: types.TypeTask, Title: "Schedule D-check for N-471SK (A330)",
		Status: types.StatusTodo, Priority: types.PriorityCritical, Parent: epic3,
		Tags: []string{"maintenance", "heavy-check"},
	}); err != nil {
		return created, err
	}
	created++

	if _, err := seedTask(root, types.TaskFrontmatter{
		Type: types.TypeTask, Title: "Replace landing gear on three 737-800s",
		Status: types.StatusTodo, Priority: types.PriorityHigh, Parent: epic3,
		Tags: []string{"maintenance", "safety"},
	}); err != nil {
		return created, err
	}
	created++

	if _, err := seedTask(root, types.TaskFrontmatter{
		Type: types.TypeTask, Title: "Engine borescope inspection — fleet-wide",
		Status: types.StatusBacklog, Priority: types.PriorityHigh, Parent: epic3,
		Tags: []string{"maintenance", "engine"},
	}); err != nil {
		return created, err
	}
	created++

	if _, err := seedTask(root, types.TaskFrontmatter{
		Type: types.TypeTask, Title: "Update avionics software to Nav DB cycle 2406",
		Status: types.StatusCancelled, Priority: types.PriorityMedium, Parent: epic3,
		Tags: []string{"avionics"},
	}); err != nil {
		return created, err
	}
	created++

	// bugs and incidents
	if _, err := seedTask(root, types.TaskFrontmatter{
		Type: types.TypeBug, Title: "Check-in kiosk freezes at bag-drop confirmation screen",
		Status: types.StatusInProgress, Priority: types.PriorityCritical,
		Tags: []string{"ground-ops", "kiosk"},
	}); err != nil {
		return created, err
	}
	created++

	if _, err := seedTask(root, types.TaskFrontmatter{
		Type: types.TypeBug, Title: "Boarding passes not scanning at gate B14",
		Status: types.StatusDone, Priority: types.PriorityCritical,
		Tags: []string{"ground-ops", "gate"},
	}); err != nil {
		return created, err
	}
	created++

	if _, err := seedTask(root, types.TaskFrontmatter{
		Type: types.TypeBug, Title: "IFE screens stuck on startup logo on rows 30–34",
		Status: types.StatusTodo, Priority: types.PriorityHigh, Parent: epic2,
		Tags: []string{"cabin", "ife"},
	}); err != nil {
		return created, err
	}
	created++

	if _, err := seedTask(root, types.TaskFrontmatter{
		Type: types.TypeBug, Title: "Meal preference not saved when booking via mobile app",
		Status: types.StatusBacklog, Priority: types.PriorityMedium,
		Tags: []string{"passenger", "catering"},
	}); err != nil {
		return created, err
	}
	created++

	return created, nil
}

// seedTask allocates an id via the project counter and writes one demo task.
func seedTask(root string, fm types.TaskFrontmatter) (string, error) {
	// allocate id + filename from the counter
	counter, err := id.NextCounter(root)
	if err != nil {
		return "", err
	}
	fm.ID = id.GenerateID(fm.Type, counter)
	fm.Created = today()
	fm.Updated = today()
	filename := id.GenerateFilename(fm.Type, counter, fm.Title)

	// write the task file
	if _, err := store.CreateTask(root, fm, filename, ""); err != nil {
		return "", err
	}
	return fm.ID, nil
}

// CreateTask creates a single task via the project counter and returns it.
// Used by the POST /api/create endpoint.
func CreateTask(root, title string, taskType types.TaskType, priority types.Priority, parent string) (*types.Task, error) {
	// build the frontmatter, optionally setting the parent
	fm := types.TaskFrontmatter{
		Type:     taskType,
		Title:    title,
		Status:   types.StatusTodo,
		Priority: priority,
	}
	if parent != "" {
		fm.Parent = parent
	}

	// create via the counter, then reload the full task
	id, err := seedTask(root, fm)
	if err != nil {
		return nil, err
	}
	summary, err := store.FindByID(root, id)
	if err != nil {
		return nil, err
	}
	task, err := store.LoadTask(summary.FilePath)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func today() string {
	return time.Now().UTC().Format("2006-01-02")
}
