package glockify

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
)

// ProjectNode manipulating Project resource.
type ProjectNode struct {
	endpoint string
	apiKey   string
}

// Project represent Clockify's project resource.
// See: https://clockify.me/developers-api#tag-Project
type Project struct {
	ID             string         `json:"id,omitempty"`
	Name           string         `json:"name,omitempty"`
	HourlyRate     HourlyRate     `json:"hourlyRate,omitempty"`
	ClientID       string         `json:"clientId,omitempty"`
	Client         string         `json:"client,omitempty"`
	WorkspaceID    string         `json:"workspaceId,omitempty"`
	Billable       bool           `json:"billable,omitempty"`
	Memberships    []Memberships  `json:"memberships,omitempty"`
	Color          string         `json:"color,omitempty"`
	Estimate       Estimate       `json:"estimate,omitempty"`
	Archived       bool           `json:"archived,omitempty"`
	Tasks          []Tasks        `json:"tasks,omitempty"`
	Note           string         `json:"note,omitempty"`
	Duration       string         `json:"duration,omitempty"`
	CostRate       int            `json:"costRate,omitempty"`
	TimeEstimate   TimeEstimate   `json:"timeEstimate,omitempty"`
	BudgetEstimate BudgetEstimate `json:"budgetEstimate"`
	CustomFields   []CustomFields `json:"customFields,omitempty"`
	Public         bool           `json:"public,omitempty"`
	Template       bool           `json:"template,omitempty"`
	Favorite       bool           `json:"favorite,omitempty"`
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

const (
	hydratedKey       = "hydrated"
	clientsKey        = "clients"
	containsClientKey = "contains-client"
	clientStatusKey   = "client-status"
	usersKey          = "users"
	containsUserKey   = "contains-users"
	userStatusKey     = "user-status"
	isTemplateKey     = "is-template"
	estimateTypeKey   = "estimate-type"
	clientIDKey       = "clientId"
	isPublicKey       = "isPublic"
	colorKey          = "color"
	noteKey           = "note"
	hourlyRateKey     = "hourlyRate"
	timeEstimateKey   = "timeEstimate"
	budgetEstimateKey = "budgetEstimate"
	membershipsKey    = "memberships"
)

type ClientStatus string

// Possible values of ClientStatus
const (
	ClientStatusActive   ClientStatus = "ACTIVE"
	ClientStatusArchived              = "ARCHIVED"
)

type UserStatus string

// Possible values of UserStatus
const (
	UserStatusActive   UserStatus = "ACTIVE"
	UserStatusInactive            = "INACTIVE"
)

type EstimateType string

// Possible values of EstimateType
const (
	EstimateTypeManual EstimateType = "MANUAL"
	EstimateTypeAuto                = "AUTO"
)

// WithHydrated if set to true, projects returned will contain custom fields,
// task and memberships. Default to false.
func WithHydrated(hydrated bool) RequestOption {
	return RequestOption{
		paramsProvider: func(v url.Values) string {
			v.Set(hydratedKey, strconv.FormatBool(hydrated))
			return hydratedKey
		},
	}
}

// WithClients if set, projects will be filtered by client IDs.
// Filter behaviour depends on the WithContainsClient.
func WithClients(ids []string) RequestOption {
	return RequestOption{
		paramsProvider: func(v url.Values) string {
			for _, id := range ids {
				v.Add(clientsKey, id)
			}
			return clientsKey
		},
	}
}

// WithContainsClient if set to true, WithClients filter will be inclusion,
// otherwise it will be exclusion. Default to true.
func WithContainsClient(containsClient bool) RequestOption {
	return RequestOption{
		paramsProvider: func(v url.Values) string {
			v.Set(containsClientKey, strconv.FormatBool(containsClient))
			return containsClientKey
		},
	}
}

// WithClientStatus filter projects returned with client state.
func WithClientStatus(status ClientStatus) RequestOption {
	return RequestOption{
		paramsProvider: func(v url.Values) string {
			v.Set(clientStatusKey, string(status))
			return clientStatusKey
		},
	}
}

// WithUsers if set, projects will be filtered by user IDs who have access.
// Filter behaviour depends on the WithContainsUser.
func WithUsers(ids []string) RequestOption {
	return RequestOption{
		paramsProvider: func(v url.Values) string {
			for _, id := range ids {
				v.Add(usersKey, id)
			}
			return usersKey
		},
	}
}

// WithContainsUser if set to true, WithUsers filter will be inclusion,
// otherwise it will be exclusion. Default to true.
func WithContainsUser(containsClient bool) RequestOption {
	return RequestOption{
		paramsProvider: func(v url.Values) string {
			v.Set(containsUserKey, strconv.FormatBool(containsClient))
			return containsUserKey
		},
	}
}

// WithUserStatus filter projects returned with user state.
func WithUserStatus(status UserStatus) RequestOption {
	return RequestOption{
		paramsProvider: func(v url.Values) string {
			v.Set(userStatusKey, string(status))
			return userStatusKey
		},
	}
}

// WithIsTemplate when applied to ProjectNode.All, it's filter projects returned
// whether it's used for template. When applied to ProjectNode.UpdateTemplate
// it's set whether project used for template.
func WithIsTemplate(isTemplate bool) RequestOption {
	return RequestOption{
		paramsProvider: func(v url.Values) string {
			v.Set(isTemplateKey, strconv.FormatBool(isTemplate))
			return isTemplateKey
		},
	}
}

// WithEstimateType
// EstimateTypeManual: type enables one fixed estimate for the whole project.
// EstimateTypeAuto: type enables task-based project estimate, and
// estimate duration doesn't matter.
func WithEstimateType(estimateType EstimateType) RequestOption {
	return RequestOption{
		paramsProvider: func(v url.Values) string {
			v.Set(estimateTypeKey, string(estimateType))
			return estimateTypeKey
		},
	}
}

// WithClientID set project's client id.
func WithClientID(id string) RequestOption {
	return RequestOption{
		paramsProvider: func(v url.Values) string {
			v.Set(clientIDKey, id)
			return clientIDKey
		},
	}
}

// WithIsPublic set project's permission state.
func WithIsPublic(isPublic bool) RequestOption {
	return RequestOption{
		paramsProvider: func(v url.Values) string {
			v.Set(isPublicKey, strconv.FormatBool(isPublic))
			return isPublicKey
		},
	}
}

// WithColor set project's color. It's in hex format, ex: #ffffff for white.
func WithColor(color string) RequestOption {
	return RequestOption{
		paramsProvider: func(v url.Values) string {
			v.Set(colorKey, color)
			return colorKey
		},
	}
}

// WithNote set project's note.
func WithNote(note string) RequestOption {
	return RequestOption{
		paramsProvider: func(v url.Values) string {
			v.Set(noteKey, note)
			return noteKey
		},
	}
}

// WithHourlyRate set project's hourly rates.
func WithHourlyRate(rate HourlyRate) RequestOption {
	return RequestOption{
		paramsProvider: func(v url.Values) string {
			hr, err := json.Marshal(rate)
			if err != nil {
				log.Fatalf("%v", err)
			}
			v.Set(hourlyRateKey, string(hr))
			return hourlyRateKey
		},
	}
}

// WithTimeEstimate set project's time estimate.
func WithTimeEstimate(estimate TimeEstimate) RequestOption {
	return RequestOption{
		paramsProvider: func(v url.Values) string {
			te, err := json.Marshal(estimate)
			if err != nil {
				log.Fatalf("%v", err)
			}
			v.Set(timeEstimateKey, string(te))
			return timeEstimateKey
		},
	}
}

// WithBudgetEstimate set project's budget estimate.
func WithBudgetEstimate(estimate BudgetEstimate) RequestOption {
	return RequestOption{
		paramsProvider: func(v url.Values) string {
			be, err := json.Marshal(estimate)
			if err != nil {
				log.Fatalf("%v", err)
			}
			v.Set(budgetEstimateKey, string(be))
			return budgetEstimateKey
		},
	}
}

// WithMemberships set project's membership state.
func WithMemberships(memberships Memberships) RequestOption {
	return RequestOption{
		paramsProvider: func(v url.Values) string {
			m, err := json.Marshal(memberships)
			if err != nil {
				log.Fatalf("%v", err)
			}
			v.Set(membershipsKey, string(m))
			return membershipsKey
		},
	}
}

type ProjectSortColumn string

const (
	ProjectSortColumnName       ProjectSortColumn = "NAME"
	ProjectSortColumnClientName                   = "CLIENT_NAME"
	ProjectSortColumnDuration                     = "DURATION"
)

// WithProjectSortColumn set fields you want to sort against.
func WithProjectSortColumn(sortColumn ProjectSortColumn) RequestOption {
	return RequestOption{
		paramsProvider: func(v url.Values) string {
			v.Set(sortColumnKey, string(sortColumn))
			return sortColumnKey
		},
	}
}

type projectAddFields struct {
	Name     string `json:"name"`
	ClientID string `json:"clientId,omitempty"`
	IsPublic *bool  `json:"isPublic,omitempty"`
	Color    string `json:"color,omitempty"`
	Note     string `json:"note,omitempty"`
	Billable *bool  `json:"billable,omitempty"`
	Public   *bool  `json:"public,omitempty"`
}

type projectUpdateFields struct {
	Name       string      `json:"name,omitempty"`
	ClientID   string      `json:"clientId,omitempty"`
	IsPublic   *bool       `json:"isPublic,omitempty"`
	HourlyRate *HourlyRate `json:"hourlyRate,omitempty"`
	Color      string      `json:"color,omitempty"`
	Note       string      `json:"note,omitempty"`
	Billable   *bool       `json:"billable,omitempty"`
	Archived   *bool       `json:"archived,omitempty"`
}

type projectUpdateEstimateFields struct {
	TimeEstimate   *TimeEstimate   `json:"timeEstimate,omitempty"`
	BudgetEstimate *BudgetEstimate `json:"budgetEstimate,omitempty"`
}

type projectUpdateMembershipsFields struct {
	Memberships *Memberships `json:"memberships,omitempty"`
}

type projectUpdateTemplateFields struct {
	IsTemplate *bool `json:"isTemplate,omitempty"`
}

// All get all Project resource based on filter given.
func (p *ProjectNode) All(workspaceID string, opts ...RequestOption) ([]Project, error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/projects", p.endpoint, workspaceID)
	res, err := get(projectAllRequest(p.apiKey, endpoint, opts))
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

func projectAllRequest(apiKey string, endpoint string, options []RequestOption) requestOptions {
	res := requestOptions{
		apiKey:   apiKey,
		endpoint: endpoint,
	}
	res.params = url.Values{}
	res.params.Add(archivedKey, strconv.FormatBool(false))
	res.params.Add(containsClientKey, strconv.FormatBool(true))
	res.params.Add(containsUserKey, strconv.FormatBool(true))
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

// Get one Project by its id.
func (p *ProjectNode) Get(workspaceID string, id string, opts ...RequestOption) (*Project, error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/projects/%s", p.endpoint, workspaceID, id)
	res, err := get(projectGetRequest(p.apiKey, endpoint, opts))
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

func projectGetRequest(apiKey string, endpoint string, options []RequestOption) requestOptions {
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
	injectContext(&res, options)

	return res
}

// Add create new Project based on fields given.
func (p *ProjectNode) Add(workspaceID string, name string, opts ...RequestOption) (*Project,
	error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/projects", p.endpoint, workspaceID)
	res, err := post(projectAddRequest(p.apiKey, endpoint, name, opts))
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

func projectAddRequest(apiKey string, endpoint string, name string,
	options []RequestOption) requestOptions {
	res := requestOptions{
		apiKey:   apiKey,
		endpoint: endpoint,
	}
	tr := true
	fields := projectAddFields{Name: name, Billable: &tr}
	for _, opt := range options {
		if opt.paramsProvider != nil {
			params := url.Values{}
			key := opt.paramsProvider(params)
			switch key {
			case clientIDKey:
				fields.ClientID = params.Get(key)
			case isPublicKey:
				val, _ := strconv.ParseBool(params.Get(isPublicKey))
				fields.IsPublic = &val
			case colorKey:
				fields.Color = params.Get(colorKey)
			case noteKey:
				fields.Note = params.Get(noteKey)
			case billableKey:
				val, _ := strconv.ParseBool(params.Get(billableKey))
				fields.Billable = &val
			}

		}
	}
	res.fields = fields
	injectContext(&res, options)

	return res
}

// Update existing Project based on fields and options given.
func (p *ProjectNode) Update(workspaceID string, id string, opts ...RequestOption) (*Project,
	error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/projects/%s", p.endpoint, workspaceID, id)
	res, err := put(projectUpdateRequest(p.apiKey, endpoint, opts))
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

func projectUpdateRequest(apiKey string, endpoint string, options []RequestOption) requestOptions {
	res := requestOptions{
		apiKey:   apiKey,
		endpoint: endpoint,
	}
	res.params = url.Values{}
	res.params.Add(estimateTypeKey, EstimateTypeAuto)
	for _, opt := range options {
		if opt.paramsProvider != nil {
			opt.paramsProvider(res.params)
		}
	}

	tr := true
	fields := projectUpdateFields{Billable: &tr}
	for _, opt := range options {
		if opt.paramsProvider != nil {
			params := url.Values{}
			key := opt.paramsProvider(params)
			switch key {
			case nameKey:
				fields.Name = params.Get(key)
			case clientIDKey:
				fields.ClientID = params.Get(key)
			case isPublicKey:
				val, _ := strconv.ParseBool(params.Get(key))
				fields.IsPublic = &val
			case hourlyRateKey:
				hr := &HourlyRate{}
				_ = json.Unmarshal([]byte(params.Get(key)), hr)
				fields.HourlyRate = hr
			case colorKey:
				fields.Color = params.Get(key)
			case noteKey:
				fields.Note = params.Get(key)
			case billableKey:
				val, _ := strconv.ParseBool(params.Get(key))
				fields.Billable = &val
			case archivedKey:
				val, _ := strconv.ParseBool(params.Get(key))
				fields.Archived = &val
			}
		}
	}
	res.params.Del(nameKey)
	res.params.Del(clientIDKey)
	res.params.Del(isPublicKey)
	res.params.Del(hourlyRateKey)
	res.params.Del(colorKey)
	res.params.Del(noteKey)
	res.params.Del(billableKey)
	res.params.Del(archivedKey)
	res.fields = fields

	injectContext(&res, options)

	return res
}

// UpdateEstimate update existing Project's estimate based on fields and options given.
func (p *ProjectNode) UpdateEstimate(workspaceID string, id string,
	opts ...RequestOption) (*Project, error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/projects/%s/estimate", p.endpoint,
		workspaceID, id)
	res, err := patch(projectUpdateEstimateRequest(p.apiKey, endpoint, opts))
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

func projectUpdateEstimateRequest(apiKey string, endpoint string,
	options []RequestOption) requestOptions {
	res := requestOptions{
		apiKey:   apiKey,
		endpoint: endpoint,
	}

	fields := projectUpdateEstimateFields{}
	for _, opt := range options {
		if opt.paramsProvider != nil {
			params := url.Values{}
			key := opt.paramsProvider(params)
			switch key {
			case timeEstimateKey:
				te := &TimeEstimate{}
				_ = json.Unmarshal([]byte(params.Get(key)), te)
				fields.TimeEstimate = te
			case budgetEstimateKey:
				be := &BudgetEstimate{}
				_ = json.Unmarshal([]byte(params.Get(key)), be)
				fields.BudgetEstimate = be
			}
		}
	}
	res.fields = fields

	injectContext(&res, options)

	return res
}

// UpdateMemberships update existing Project's memberships based on fields given.
func (p *ProjectNode) UpdateMemberships(workspaceID string, id string,
	opts ...RequestOption) (*Project, error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/projects/%s/memberships", p.endpoint,
		workspaceID, id)
	res, err := patch(projectUpdateMembershipRequest(p.apiKey, endpoint, opts))
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

func projectUpdateMembershipRequest(apiKey string, endpoint string,
	options []RequestOption) requestOptions {
	res := requestOptions{
		apiKey:   apiKey,
		endpoint: endpoint,
	}

	fields := projectUpdateMembershipsFields{}
	for _, opt := range options {
		if opt.paramsProvider != nil {
			params := url.Values{}
			key := opt.paramsProvider(params)
			switch key {
			case membershipsKey:
				m := &Memberships{}
				_ = json.Unmarshal([]byte(params.Get(key)), m)
				fields.Memberships = m
			}
		}
	}
	res.fields = fields

	injectContext(&res, options)

	return res
}

// UpdateTemplate update existing Project's template based on fields options given.
func (p *ProjectNode) UpdateTemplate(workspaceID string, id string,
	opts ...RequestOption) (*Project, error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/projects/%s/template", p.endpoint,
		workspaceID, id)
	res, err := patch(projectUpdateTemplateRequest(p.apiKey, endpoint, opts))
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

func projectUpdateTemplateRequest(apiKey string, endpoint string,
	options []RequestOption) requestOptions {
	res := requestOptions{
		apiKey:   apiKey,
		endpoint: endpoint,
	}

	fields := projectUpdateTemplateFields{}
	for _, opt := range options {
		if opt.paramsProvider != nil {
			params := url.Values{}
			key := opt.paramsProvider(params)
			switch key {
			case isTemplateKey:
				val, _ := strconv.ParseBool(params.Get(key))
				fields.IsTemplate = &val
			}
		}
	}
	res.fields = fields

	injectContext(&res, options)

	return res
}

// Delete existing Project.
func (p *ProjectNode) Delete(workspaceID string, id string, opts ...RequestOption) (*Project,
	error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/projects/%s", p.endpoint, workspaceID, id)
	res, err := del(projectDeleteRequest(p.apiKey, endpoint, opts))
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

func projectDeleteRequest(apiKey string, endpoint string, options []RequestOption) requestOptions {
	res := requestOptions{
		apiKey:   apiKey,
		endpoint: endpoint,
	}
	injectContext(&res, options)

	return res
}
