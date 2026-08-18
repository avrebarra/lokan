package server

import (
	"strconv"
	"time"

	"github.com/avrebarra/lokan/internal/store"
	"github.com/avrebarra/lokan/internal/types"
)

// SeedDemoData creates the demo airline tasks and returns how many tasks were
// created. IDs come from the project counter.
func SeedDemoData(board string) (int, error) {
	demo := []types.TaskFrontmatter{
		{Title: "Submit route approval to FAA and JCAB", Status: types.StatusDone, Tags: []string{"compliance", "international"}},
		{Title: "Assign Boeing 787-9 fleet for long-haul", Status: types.StatusInProgress, Tags: []string{"fleet", "operations"}},
		{Title: "Hire and train Japan-route cabin crew", Status: types.StatusTodo, Tags: []string{"crew", "training"}},
		{Title: "Launch marketing campaign for SFO–NRT", Status: types.StatusBacklog, Tags: []string{"marketing"}},
		{Title: "Redesign business class seating", Status: types.StatusDone, Tags: []string{"cabin", "design"}},
		{Title: "Roll out in-flight Wi-Fi on all A320s", Status: types.StatusInProgress, Tags: []string{"connectivity", "fleet"}},
		{Title: "Introduce premium meal service with local chefs", Status: types.StatusTodo, Tags: []string{"catering", "passenger"}},
		{Title: "Refresh loyalty program tier benefits", Status: types.StatusBacklog, Tags: []string{"loyalty", "passenger"}},
		{Title: "Schedule D-check for N-471SK (A330)", Status: types.StatusTodo, Tags: []string{"maintenance", "heavy-check"}},
		{Title: "Replace landing gear on three 737-800s", Status: types.StatusTodo, Tags: []string{"maintenance", "safety"}},
		{Title: "Check-in kiosk freezes at bag-drop confirmation screen", Status: types.StatusInProgress, Tags: []string{"ground-ops", "kiosk"}},
		{Title: "Boarding passes not scanning at gate B14", Status: types.StatusDone, Tags: []string{"ground-ops", "gate"}},
	}

	created := 0
	for _, fm := range demo {
		// allocate id from the counter
		counter, err := store.NextCounter(board)
		if err != nil {
			return created, err
		}
		fm.ID = strconv.Itoa(counter)
		fm.Created = today()
		fm.Updated = today()
		if _, err := store.CreateTask(board, fm, ""); err != nil {
			return created, err
		}
		created++
	}
	return created, nil
}

func today() string {
	return time.Now().UTC().Format("2006-01-02")
}
