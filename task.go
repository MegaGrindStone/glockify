package glockify

import (
	"context"
	"encoding/json"
	"fmt"
)

// TaskNode manipulating Task resource.
type TaskNode struct {
	workspaceID  string
	projectID    string
	baseEndpoint string
	apiKey       string
}

// Task wraps Clockify's task resource.
// See: https://clockify.me/developers-api#tag-Task
type Task struct {
	AssigneeIds []string   `json:"assigneeIds,omitempty"`
	Estimate    string     `json:"estimate,omitempty"`
	ID          string     `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`
	ProjectID   string     `json:"projectId,omitempty"`
	Billable    string     `json:"billable,omitempty"`
	HourlyRate  HourlyRate `json:"hourlyRate,omitempty"`
	CostRate    CostRate   `json:"costRate,omitempty"`
	Status      string     `json:"status,omitempty"`
}

// CostRate wraps Clockify's cost rate resource.
// See: https://clockify.me/developers-api#tag-Cost
type CostRate struct {
	Amount   string `json:"amount,omitempty"`
	Currency string `json:"currency,omitempty"`
}

// TaskAllFilter is used for All request.
type TaskAllFilter struct {
	// IsActive if provided and true, only active tasks will be returned.
	// Otherwise, only finished tasks will be returned.
	IsActive bool `schema:"is-active"`

	// Name if provided, task will be filtered by name
	Name string `schema:"name"`

	// Page default 1
	Page int `schema:"page"`

	// PageSize max page-size 5000
	// Default 50
	PageSize int `schema:"page-size"`
}

// TaskAddFields is used for Add request.
// See: https://clockify.me/developers-api#tag-Task
type TaskAddFields struct {
	Name        string   `json:"name,omitempty"`
	AssigneeIds []string `json:"assigneeIds,omitempty"`
	Estimate    string   `json:"estimate,omitempty"`
	Status      string   `json:"status,omitempty"`
}

// TaskUpdateFields is used for Update request.
// See: https://clockify.me/developers-api#tag-Task
type TaskUpdateFields struct {
	Name        string   `json:"name"`
	AssigneeIds []string `json:"assigneeIds,omitempty"`
	Estimate    string   `json:"estimate,omitempty"`
	Billable    bool     `json:"billable,omitempty"`
	Status      string   `json:"status,omitempty"`
}

// All get all Task resource based on filter given.
func (t *TaskNode) All(ctx context.Context, filter TaskAllFilter) ([]Task, error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/projects/%s/tasks", t.baseEndpoint,
		t.workspaceID, t.projectID)
	res, err := get(ctx, t.apiKey, filter, endpoint)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}
	result := make([]Task, 0)
	if err := json.Unmarshal(res, &result); err != nil {
		if jErr, ok := err.(*json.UnmarshalTypeError); ok {
			return nil, fmt.Errorf("unmarshal field %v of type %v", jErr.Field, jErr.Type)
		}
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}
	return result, nil
}

// Get one Task by its id.
func (t *TaskNode) Get(ctx context.Context, id string) (*Task, error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/projects/%s/tasks/%s", t.baseEndpoint, t.workspaceID,
		t.projectID, id)
	res, err := get(ctx, t.apiKey, nil, endpoint)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}
	result := new(Task)
	if err := json.Unmarshal(res, &result); err != nil {
		if jErr, ok := err.(*json.UnmarshalTypeError); ok {
			return nil, fmt.Errorf("unmarshal field %v of type %v", jErr.Field, jErr.Type)
		}
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}
	return result, nil
}

// Add create new Task based on fields given.
func (t *TaskNode) Add(ctx context.Context, fields TaskAddFields) (*Task, error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/projects/%s/tasks", t.baseEndpoint,
		t.workspaceID, t.projectID)
	res, err := post(ctx, t.apiKey, nil, fields, endpoint)
	if err != nil {
		return nil, fmt.Errorf("post: %w", err)
	}
	result := new(Task)
	if err := json.Unmarshal(res, &result); err != nil {
		if jErr, ok := err.(*json.UnmarshalTypeError); ok {
			return nil, fmt.Errorf("unmarshal field %v of type %v", jErr.Field, jErr.Type)
		}
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}
	return result, nil
}

// Update existing Task based on fields and options given.
func (t *TaskNode) Update(ctx context.Context, id string, fields TaskUpdateFields) (*Task, error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/projects/%s/tasks/%s", t.baseEndpoint,
		t.workspaceID, t.projectID, id)
	res, err := put(ctx, t.apiKey, nil, fields, endpoint)
	if err != nil {
		return nil, fmt.Errorf("put: %w", err)
	}
	result := new(Task)
	if err := json.Unmarshal(res, &result); err != nil {
		if jErr, ok := err.(*json.UnmarshalTypeError); ok {
			return nil, fmt.Errorf("unmarshal field %v of type %v", jErr.Field, jErr.Type)
		}
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}
	return result, nil
}

// Delete existing Task.
func (t *TaskNode) Delete(ctx context.Context, id string) (*Task, error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/projects/%s/task/%s", t.baseEndpoint,
		t.workspaceID, t.projectID, id)
	res, err := del(ctx, t.apiKey, endpoint)
	if err != nil {
		return nil, fmt.Errorf("del: %w", err)
	}
	result := new(Task)
	if err := json.Unmarshal(res, &result); err != nil {
		if jErr, ok := err.(*json.UnmarshalTypeError); ok {
			return nil, fmt.Errorf("unmarshal field %v of type %v", jErr.Field, jErr.Type)
		}
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}
	return result, nil
}
