package glockify

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// TaskNode manipulating Task resource.
type TaskNode struct {
	endpoint string
	apiKey   string
}

// Task represents Clockify's task resource.
// See: https://clockify.me/developers-api#tag-Task
type Task struct {
	AssigneeIds []string   `json:"assigneeIds,omitempty"`
	Estimate    string     `json:"estimate,omitempty"`
	ID          string     `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`
	ProjectID   string     `json:"projectId,omitempty"`
	Billable    bool       `json:"billable,omitempty"`
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

const (
	isActiveKey         = "is-active"
	strictNameSearchKey = "strict-name-search"
	assigneeIDsKey      = "assigneeIds"
	estimateKey         = "estimate"
	statusKey           = "status"
)

// WithIsActive filter task by active state.
func WithIsActive(isActive bool) RequestOption {
	return RequestOption{
		paramsProvider: func(v url.Values) string {
			v.Set(isActiveKey, strconv.FormatBool(isActive))
			return isActiveKey
		},
	}
}

// WithStrictNameSearch if set to true, WithName filter will be exact match,
// meanwhile if set to false, partial search is executed.
func WithStrictNameSearch(on bool) RequestOption {
	return RequestOption{
		paramsProvider: func(v url.Values) string {
			v.Set(strictNameSearchKey, strconv.FormatBool(on))
			return strictNameSearchKey
		},
	}
}

type TaskSortColumn string

// Possible value of TaskSortColumn
const (
	TaskSortColumnID   TaskSortColumn = "ID"
	TaskSortColumnName                = "NAME"
)

// WithTaskSortColumn set fields you want to sort against.
func WithTaskSortColumn(sortColumn TaskSortColumn) RequestOption {
	return RequestOption{
		paramsProvider: func(v url.Values) string {
			v.Set(sortColumnKey, string(sortColumn))
			return sortColumnKey
		},
	}
}

// WithAssigneeIDs set assignees for this task.
func WithAssigneeIDs(ids []string) RequestOption {
	return RequestOption{
		paramsProvider: func(v url.Values) string {
			v.Set(assigneeIDsKey, strings.Join(ids, arraySeparator))
			return assigneeIDsKey
		},
	}
}

// WithEstimate set task estimate in Clockify time format. Eg: "PT2H" for 2 hour.
func WithEstimate(estimate string) RequestOption {
	return RequestOption{
		paramsProvider: func(v url.Values) string {
			v.Set(estimateKey, estimate)
			return estimateKey
		},
	}
}

type TaskStatus string

const (
	TaskStatusActive TaskStatus = "ACTIVE"
	TaskStatusDone              = "DONE"
)

// WithStatus set task state.
func WithStatus(status TaskStatus) RequestOption {
	return RequestOption{
		paramsProvider: func(v url.Values) string {
			v.Set(statusKey, string(status))
			return statusKey
		},
	}
}

type taskAddFields struct {
	Name        string   `json:"name"`
	AssigneeIds []string `json:"assigneeIds,omitempty"`
	Estimate    string   `json:"estimate,omitempty"`
	Status      string   `json:"status,omitempty"`
}

type taskUpdateFields struct {
	Name        string   `json:"name,omitempty"`
	AssigneeIds []string `json:"assigneeIds,omitempty"`
	Estimate    string   `json:"estimate,omitempty"`
	Billable    *bool    `json:"billable,omitempty"`
	Status      string   `json:"status,omitempty"`
}

// All get all Task resource based on filter given.
func (t *TaskNode) All(workspaceID string, projectID string, opts ...RequestOption) ([]Task,
	error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/projects/%s/tasks", t.endpoint,
		workspaceID, projectID)
	res, err := get(taskAllRequest(t.apiKey, endpoint, opts))
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

func taskAllRequest(apiKey string, endpoint string, options []RequestOption) requestOptions {
	res := requestOptions{
		apiKey:   apiKey,
		endpoint: endpoint,
	}
	res.params = url.Values{}
	res.params.Add(pageKey, strconv.Itoa(defaultPage))
	res.params.Add(pageSizeKey, strconv.Itoa(defaultPageSize))
	res.params.Add(sortOrderKey, defaultSortOrder)
	for _, opt := range options {
		if opt.paramsProvider != nil {
			opt.paramsProvider(res.params)
		}
	}
	injectContext(&res, options)

	return res
}

// Get one Task by its id.
func (t *TaskNode) Get(workspaceID string, projectID string, id string,
	opts ...RequestOption) (*Task, error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/projects/%s/tasks/%s", t.endpoint, workspaceID,
		projectID, id)
	res, err := get(taskGetRequest(t.apiKey, endpoint, opts))
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

func taskGetRequest(apiKey string, endpoint string, options []RequestOption) requestOptions {
	res := requestOptions{
		apiKey:   apiKey,
		endpoint: endpoint,
	}
	res.params = url.Values{}
	injectContext(&res, options)

	return res
}

// Add create new Task based on fields given.
func (t *TaskNode) Add(workspaceID string, projectID string, name string,
	opts ...RequestOption) (*Task, error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/projects/%s/tasks", t.endpoint,
		workspaceID, projectID)
	res, err := post(taskAddRequest(t.apiKey, endpoint, name, opts))
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

func taskAddRequest(apiKey string, endpoint string, name string,
	options []RequestOption) requestOptions {
	res := requestOptions{
		apiKey:   apiKey,
		endpoint: endpoint,
	}
	fields := taskAddFields{Name: name}
	for _, opt := range options {
		if opt.paramsProvider != nil {
			params := url.Values{}
			key := opt.paramsProvider(params)
			switch key {
			case assigneeIDsKey:
				fields.AssigneeIds = strings.Split(params.Get(key), arraySeparator)
			case estimateKey:
				fields.Estimate = params.Get(key)
			case statusKey:
				fields.Status = params.Get(key)
			}

		}
	}
	res.fields = fields
	injectContext(&res, options)

	return res
}

// Update existing Task based on fields and options given.
func (t *TaskNode) Update(workspaceID string, projectID string, id string,
	opts ...RequestOption) (*Task, error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/projects/%s/tasks/%s", t.endpoint,
		workspaceID, projectID, id)
	res, err := put(taskUpdateRequest(t.apiKey, endpoint, opts))
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

func taskUpdateRequest(apiKey string, endpoint string, options []RequestOption) requestOptions {
	res := requestOptions{
		apiKey:   apiKey,
		endpoint: endpoint,
	}
	res.params = url.Values{}
	for _, opt := range options {
		if opt.paramsProvider != nil {
			opt.paramsProvider(res.params)
		}
	}

	fields := taskUpdateFields{}
	for _, opt := range options {
		if opt.paramsProvider != nil {
			params := url.Values{}
			key := opt.paramsProvider(params)
			switch key {
			case nameKey:
				fields.Name = params.Get(key)
			case assigneeIDsKey:
				fields.AssigneeIds = strings.Split(params.Get(key), arraySeparator)
			case estimateKey:
				fields.Estimate = params.Get(key)
			case billableKey:
				val, _ := strconv.ParseBool(params.Get(key))
				fields.Billable = &val
			case statusKey:
				fields.Status = params.Get(key)
			}

		}
	}
	res.params.Del(nameKey)
	res.params.Del(assigneeIDsKey)
	res.params.Del(estimateKey)
	res.params.Del(billableKey)
	res.params.Del(statusKey)
	res.fields = fields

	injectContext(&res, options)

	return res
}

// Delete existing Task.
func (t *TaskNode) Delete(workspaceID string, projectID string, id string,
	opts ...RequestOption) (*Task, error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/projects/%s/task/%s", t.endpoint,
		workspaceID, projectID, id)
	res, err := del(taskDeleteRequest(t.apiKey, endpoint, opts))
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

func taskDeleteRequest(apiKey string, endpoint string, options []RequestOption) requestOptions {
	res := requestOptions{
		apiKey:   apiKey,
		endpoint: endpoint,
	}
	injectContext(&res, options)

	return res
}
