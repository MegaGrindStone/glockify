package glockify

import (
	"context"
	"encoding/json"
	"fmt"
)

// ProjectNode manipulating Project resource.
type ProjectNode struct {
	workspaceID  string
	baseEndpoint string
	apiKey       string
}

// Project wraps Clockify's project resource.
// See: https://clockify.me/developers-api#tag-Project
type Project []struct {
	ID           string        `json:"id,omitempty"`
	Name         string        `json:"name,omitempty"`
	HourlyRate   HourlyRate    `json:"hourlyRate,omitempty"`
	ClientID     string        `json:"clientId,omitempty"`
	Client       string        `json:"client,omitempty"`
	WorkspaceID  string        `json:"workspaceId,omitempty"`
	Billable     bool          `json:"billable,omitempty"`
	Memberships  []Memberships `json:"memberships,omitempty"`
	Color        string        `json:"color,omitempty"`
	Estimate     Estimate      `json:"estimate,omitempty"`
	Archived     bool          `json:"archived,omitempty"`
	Tasks        []Tasks       `json:"tasks,omitempty"`
	Note         string        `json:"note,omitempty"`
	Duration     string        `json:"duration,omitempty"`
	CostRate     int           `json:"costRate,omitempty"`
	TimeEstimate TimeEstimate  `json:"timeEstimate,omitempty"`
	//BudgetEstimate interface{}    `json:"budgetEstimate"`
	CustomFields []CustomFields `json:"customFields,omitempty"`
	Public       bool           `json:"public,omitempty"`
	Template     bool           `json:"template,omitempty"`
	Favorite     bool           `json:"favorite,omitempty"`
}

// Estimate wraps Clockify's estimate resource.
// See: https://clockify.me/developers-api#tag-Project
type Estimate struct {
	Estimate string `json:"estimate,omitempty"`
	Type     string `json:"type,omitempty"`
}

// Tasks wraps Clockify's tasks resource.
// See: https://clockify.me/developers-api#tag-Project
type Tasks struct {
	ID           string   `json:"id,omitempty"`
	Name         string   `json:"name,omitempty"`
	ProjectID    string   `json:"projectId,omitempty"`
	AssigneeIds  []string `json:"assigneeIds,omitempty"`
	AssigneeID   string   `json:"assigneeId,omitempty"`
	UserGroupIds []string `json:"userGroupIds,omitempty"`
	Estimate     string   `json:"estimate,omitempty"`
	Status       string   `json:"status,omitempty"`
	Duration     int      `json:"duration,omitempty"`
	Billable     bool     `json:"billable,omitempty"`
	HourlyRate   int      `json:"hourlyRate,omitempty"`
	CostRate     int      `json:"costRate,omitempty"`
}

// TimeEstimate wraps Clockify's time estimate resource.
// See: https://clockify.me/developers-api#tag-Project
type TimeEstimate struct {
	Estimate           string `json:"estimate,omitempty"`
	Type               string `json:"type,omitempty"`
	ResetOption        string `json:"resetOption,omitempty"`
	Active             bool   `json:"active,omitempty"`
	IncludeNonBillable bool   `json:"includeNonBillable,omitempty"`
}

// CustomFields wraps Clockify's custom fields resource.
// See: https://clockify.me/developers-api#tag-Project
type CustomFields struct {
	CustomFieldID string      `json:"customFieldId,omitempty"`
	Name          string      `json:"name,omitempty"`
	Type          string      `json:"type,omitempty"`
	Value         interface{} `json:"value,omitempty"`
	Status        string      `json:"status,omitempty"`
}

// BudgetEstimate wraps Clockify's budget estimate resource.
// See: https://clockify.me/developers-api#tag-Project
type BudgetEstimate struct {
	Estimate    string `json:"estimate,omitempty"`
	Type        string `json:"type,omitempty"`
	ResetOption string `json:"resetOption,omitempty"`
	Active      bool   `json:"active,omitempty"`
}

// ProjectAllFilter is used for All request.
type ProjectAllFilter struct {
	// Hydrated If true, you'll get custom fields, tasks, and memberships of all projects.
	Hydrated bool `schema:"hydrated"`

	// Archived if true, you'll get only archived clients.
	// If false, you'll get only active clients.
	Archived bool `schema:"archived"`

	// Name if provided, clients will be filtered by name
	Name string `schema:"name"`

	// Page default 1
	Page int `schema:"page"`

	// PageSize max page-size 5000
	// Default 50
	PageSize int `schema:"page-size"`

	// Billable If provided, projects will be filtered by billable status.
	Billable bool `schema:"billable"`

	// Clients If provided, projects will be filtered by client ID(s).
	Clients []string `schema:"clients"`

	// ContainsClient Default is true. If set to false, you'll get projects that
	// exclude the client ID(s) you've provided.
	ContainsClient bool `schema:"contains-client"`

	// ClientStatus possible value is "ACTIVE", "ARCHIVED"
	ClientStatus string `schema:"client-status"`

	// Users If provided, projects will be filtered by user ID(s) who have access.
	Users []string `schema:"users"`

	// ContainUsers Default is true. If set to false, you'll get projects that
	// the user ID(s) you've provided don't have access to.
	ContainUsers bool `schema:"contain-users"`

	// UserStatus possible value is "ACTIVE", "INACTIVE"
	UserStatus string `schema:"user-status"`

	// IsTemplate If provided, projects will be filtered by whether they are used as a template.
	IsTemplate bool `schema:"is-template"`

	// SortColumn possible values is "NAME"
	SortColumn string `schema:"sort-column"`

	// SortOrder possible values is "ASCENDING", "DESCENDING"
	SortOrder string `schema:"sort-order"`
}

// ProjectAddFields is used for Add request.
// See: https://clockify.me/developers-api#tag-Project
type ProjectAddFields struct {
	Name     string `json:"name"`
	ClientID string `json:"clientId,omitempty"`
	IsPublic string `json:"isPublic,omitempty"`
	Color    string `json:"color,omitempty"`
	Note     string `json:"note,omitempty"`
	Billable bool   `json:"billable,omitempty"`
	Public   bool   `json:"public,omitempty"`
}

// ProjectUpdateFields is used for Update request.
// See: https://clockify.me/developers-api#tag-Project
type ProjectUpdateFields struct {
	Name       string     `json:"name,omitempty"`
	ClientID   string     `json:"clientId,omitempty"`
	IsPublic   bool       `json:"isPublic,omitempty"`
	HourlyRate HourlyRate `json:"hourlyRate"`
	Color      string     `json:"color,omitempty"`
	Note       string     `json:"note,omitempty"`
	Billable   bool       `json:"billable,omitempty"`
	Archived   bool       `json:"archived,omitempty"`
}

// ProjectUpdateEstimateFields is used for UpdateEstimate request.
// See: https://clockify.me/developers-api#tag-Project
type ProjectUpdateEstimateFields struct {
	TimeEstimate   TimeEstimate   `json:"timeEstimate"`
	BudgetEstimate BudgetEstimate `json:"budgetEstimate"`
}

// ProjectUpdateMembershipsFields is used for UpdateMemberships request.
// See: https://clockify.me/developers-api#tag-Project
type ProjectUpdateMembershipsFields struct {
	Memberships Memberships `json:"memberships"`
}

// ProjectUpdateTemplateFields is used for UpdateTemplate request.
// See: https://clockify.me/developers-api#tag-Project
type ProjectUpdateTemplateFields struct {
	IsTemplate bool `json:"isTemplate"`
}

// ProjectUpdateOptions is used for UpdateEstimate request.
type ProjectUpdateOptions struct {
	// EstimateType possible values is:
	// "MANUAL": type enables one fixed estimate for the whole project.
	// "AUTO": type enables task-based project estimate.
	// If AUTO is enabled, estimate duration doesn't matter.
	EstimateType string `scheme:"estimate-type"`
}

// ProjectUpdateEstimateOptions is used for Update request.
type ProjectUpdateEstimateOptions struct {
	// Active possible values is: "time" and "budget".
	// If you need "No estimate", then don't set this field,
	// or set both ProjectUpdateFields active fields as false.
	Active string `scheme:"active"`

	// Reset possible value is "MONTHLY".
	Reset string `schema:"reset"`

	// Type possible value is:
	// "MANUAL": estimating whole project.
	// "AUTO": enable task-based estimate
	Type string `scheme:"type"`
}

// All get all Project resource based on filter given.
func (p *ProjectNode) All(ctx context.Context, filter ProjectAllFilter) ([]Project, error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/projects", p.baseEndpoint, p.workspaceID)
	res, err := get(ctx, p.apiKey, filter, endpoint)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}
	result := make([]Project, 0)
	if err := json.Unmarshal(res, &result); err != nil {
		if jErr, ok := err.(*json.UnmarshalTypeError); ok {
			return nil, fmt.Errorf("unmarshal field %v of type %v", jErr.Field, jErr.Type)
		}
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}
	return result, nil
}

// Get one Project by its id.
func (p *ProjectNode) Get(ctx context.Context, id string) (*Project, error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/projects/%s", p.baseEndpoint, p.workspaceID, id)
	res, err := get(ctx, p.apiKey, nil, endpoint)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}
	result := new(Project)
	if err := json.Unmarshal(res, &result); err != nil {
		if jErr, ok := err.(*json.UnmarshalTypeError); ok {
			return nil, fmt.Errorf("unmarshal field %v of type %v", jErr.Field, jErr.Type)
		}
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}
	return result, nil
}

// Add create new Project based on fields given.
func (p *ProjectNode) Add(ctx context.Context, fields ProjectAddFields) (*Project, error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/projects", p.baseEndpoint, p.workspaceID)
	res, err := post(ctx, p.apiKey, nil, fields, endpoint)
	if err != nil {
		return nil, fmt.Errorf("post: %w", err)
	}
	result := new(Project)
	if err := json.Unmarshal(res, &result); err != nil {
		if jErr, ok := err.(*json.UnmarshalTypeError); ok {
			return nil, fmt.Errorf("unmarshal field %v of type %v", jErr.Field, jErr.Type)
		}
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}
	return result, nil
}

// Update existing Project based on fields and options given.
func (p *ProjectNode) Update(ctx context.Context, id string, fields ProjectUpdateFields,
	options ProjectUpdateOptions) (*Project, error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/projects/%s", p.baseEndpoint, p.workspaceID, id)
	res, err := put(ctx, p.apiKey, options, fields, endpoint)
	if err != nil {
		return nil, fmt.Errorf("put: %w", err)
	}
	result := new(Project)
	if err := json.Unmarshal(res, &result); err != nil {
		if jErr, ok := err.(*json.UnmarshalTypeError); ok {
			return nil, fmt.Errorf("unmarshal field %v of type %v", jErr.Field, jErr.Type)
		}
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}
	return result, nil
}

// UpdateEstimate update existing Project's estimate based on fields and options given.
func (p *ProjectNode) UpdateEstimate(ctx context.Context, id string,
	fields ProjectUpdateEstimateFields, options ProjectUpdateEstimateOptions) (*Project, error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/projects/%s/estimate", p.baseEndpoint,
		p.workspaceID, id)
	res, err := patch(ctx, p.apiKey, options, fields, endpoint)
	if err != nil {
		return nil, fmt.Errorf("put: %w", err)
	}
	result := new(Project)
	if err := json.Unmarshal(res, &result); err != nil {
		if jErr, ok := err.(*json.UnmarshalTypeError); ok {
			return nil, fmt.Errorf("unmarshal field %v of type %v", jErr.Field, jErr.Type)
		}
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}
	return result, nil
}

// UpdateMemberships update existing Project's memberships based on fields given.
func (p *ProjectNode) UpdateMemberships(ctx context.Context, id string,
	fields ProjectUpdateMembershipsFields) (*Project, error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/projects/%s/memberships", p.baseEndpoint,
		p.workspaceID, id)
	res, err := patch(ctx, p.apiKey, nil, fields, endpoint)
	if err != nil {
		return nil, fmt.Errorf("put: %w", err)
	}
	result := new(Project)
	if err := json.Unmarshal(res, &result); err != nil {
		if jErr, ok := err.(*json.UnmarshalTypeError); ok {
			return nil, fmt.Errorf("unmarshal field %v of type %v", jErr.Field, jErr.Type)
		}
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}
	return result, nil
}

// UpdateTemplate update existing Project's template based on fields options given.
func (p *ProjectNode) UpdateTemplate(ctx context.Context, id string,
	fields ProjectUpdateTemplateFields) (*Project, error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/projects/%s/template", p.baseEndpoint,
		p.workspaceID, id)
	res, err := patch(ctx, p.apiKey, nil, fields, endpoint)
	if err != nil {
		return nil, fmt.Errorf("put: %w", err)
	}
	result := new(Project)
	if err := json.Unmarshal(res, &result); err != nil {
		if jErr, ok := err.(*json.UnmarshalTypeError); ok {
			return nil, fmt.Errorf("unmarshal field %v of type %v", jErr.Field, jErr.Type)
		}
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}
	return result, nil
}

// Delete existing Project.
func (p *ProjectNode) Delete(ctx context.Context, id string) (*Project, error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/projects/%s", p.baseEndpoint, p.workspaceID, id)
	res, err := del(ctx, p.apiKey, endpoint)
	if err != nil {
		return nil, fmt.Errorf("del: %w", err)
	}
	result := new(Project)
	if err := json.Unmarshal(res, &result); err != nil {
		if jErr, ok := err.(*json.UnmarshalTypeError); ok {
			return nil, fmt.Errorf("unmarshal field %v of type %v", jErr.Field, jErr.Type)
		}
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}
	return result, nil
}
