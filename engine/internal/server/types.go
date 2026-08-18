package server

import "github.com/avrebarra/lokan/internal/types"

// Request types — payload contracts for the write endpoints.
type RequestDataCreate struct {
	Title    string         `json:"title"`
	Type     types.TaskType `json:"type"`
	Priority types.Priority `json:"priority"`
	Parent   string         `json:"parent"`
}

type RequestDataUpdate struct {
	ID    string `json:"id"`
	Field string `json:"field"`
	Value string `json:"value"`
}

type RequestDataMove struct {
	ID       string       `json:"id"`
	Status   types.Status `json:"status"`
	BeforeID string       `json:"beforeId"`
}

type RequestDataConfigStatuses struct {
	Statuses []types.StatusDef `json:"statuses"`
}

type RequestDataClear struct {
	Scope string `json:"scope"`
}

// Response types — payload contracts for the API responses.
type ResponseDataTasks struct {
	Tasks     []types.TaskSummary `json:"tasks"`
	Statuses  []types.StatusDef   `json:"statuses"`
	Root      string              `json:"root"`
	BoardPath string              `json:"board_path"`
	BoardRoot string              `json:"board_root"`
}

type ResponseDataTask struct {
	Task types.Task `json:"task"`
}

type ResponseDataCreate struct {
	Task types.Task `json:"task"`
}

type ResponseDataUpdate struct {
	Task types.Task `json:"task"`
}

type ResponseDataMove struct {
	Task types.Task `json:"task"`
}

type ResponseDataConfigStatuses struct {
	Statuses []types.StatusDef `json:"statuses"`
	Moved    int               `json:"moved"`
}

type ResponseDataClear struct {
	Deleted int `json:"deleted"`
}

type ResponseDataSeed struct {
	Created int `json:"created"`
}

type ResponseDataError struct {
	Error string `json:"error"`
}
