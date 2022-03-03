package glockify

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// ClientNode manipulating Client resource.
type ClientNode struct {
	endpoint string
	apiKey   string
}

// Client represent Clockify's client resource.
// See: https://clockify.me/developers-api#tag-Client
type Client struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	WorkspaceID string `json:"workspaceId,omitempty"`
	Archived    bool   `json:"archived,omitempty"`
}

const (
	archiveProjectsKey = "archive-projects"
)

// WithArchiveProjects set whether archiving client will result in archiving
// all projects of given client. Default to false.
func WithArchiveProjects(archiveProjects bool) RequestOption {
	return RequestOption{
		paramsProvider: func(v url.Values) string {
			v.Set(archiveProjectsKey, strconv.FormatBool(archiveProjects))
			return archiveProjectsKey
		},
	}
}

type ClientSortColumn string

const (
	ClientSortColumnName ClientSortColumn = "NAME"
)

// WithClientSortColumn set fields you want to sort against.
func WithClientSortColumn(sortColumn ClientSortColumn) RequestOption {
	return RequestOption{
		paramsProvider: func(v url.Values) string {
			v.Set(sortColumnKey, string(sortColumn))
			return sortColumnKey
		},
	}
}

type clientAddFields struct {
	Name string `json:"name"`
}

type clientUpdateFields struct {
	Archived *bool  `json:"archived,omitempty"`
	Name     string `json:"name,omitempty"`
}

// All get all Client resource based on filter given.
func (c *ClientNode) All(workspaceID string, opts ...RequestOption) ([]Client, error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/clients", c.endpoint, workspaceID)
	res, err := get(clientAllRequest(c.apiKey, endpoint, opts))
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}
	result := make([]Client, 0)
	if err := json.Unmarshal(res, &result); err != nil {
		if jErr, ok := err.(*json.UnmarshalTypeError); ok {
			return nil, fmt.Errorf("unmarshal field %v of type %v", jErr.Field, jErr.Type)
		}
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}
	return result, nil
}

func clientAllRequest(apiKey string, endpoint string, options []RequestOption) requestOptions {
	res := requestOptions{
		apiKey:   apiKey,
		endpoint: endpoint,
	}
	res.params = url.Values{}
	res.params.Add(archivedKey, strconv.FormatBool(false))
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

// Get one Client by its id.
func (c *ClientNode) Get(workspaceID string, id string, opts ...RequestOption) (*Client, error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/clients/%s", c.endpoint, workspaceID, id)
	res, err := get(clientGetRequest(c.apiKey, endpoint, opts))
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}
	result := new(Client)
	if err := json.Unmarshal(res, &result); err != nil {
		if jErr, ok := err.(*json.UnmarshalTypeError); ok {
			return nil, fmt.Errorf("unmarshal field %v of type %v", jErr.Field, jErr.Type)
		}
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}
	return result, nil
}

func clientGetRequest(apiKey string, endpoint string, options []RequestOption) requestOptions {
	res := requestOptions{
		apiKey:   apiKey,
		endpoint: endpoint,
	}
	res.params = url.Values{}
	injectContext(&res, options)

	return res
}

// Add create new Client based on fields given.
func (c *ClientNode) Add(workspaceID string, name string, opts ...RequestOption) (*Client, error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/clients", c.endpoint, workspaceID)
	res, err := post(clientAddRequest(c.apiKey, endpoint, name, opts))
	if err != nil {
		return nil, fmt.Errorf("post: %w", err)
	}
	result := new(Client)
	if err := json.Unmarshal(res, &result); err != nil {
		if jErr, ok := err.(*json.UnmarshalTypeError); ok {
			return nil, fmt.Errorf("unmarshal field %v of type %v", jErr.Field, jErr.Type)
		}
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}
	return result, nil
}

func clientAddRequest(apiKey string, endpoint string, name string,
	options []RequestOption) requestOptions {
	res := requestOptions{
		apiKey:   apiKey,
		endpoint: endpoint,
	}
	res.fields = clientAddFields{Name: name}
	injectContext(&res, options)

	return res
}

// Update existing Client based on fields and options given.
func (c *ClientNode) Update(workspaceID string, id string, opts ...RequestOption) (*Client,
	error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/clients/%s", c.endpoint, workspaceID, id)
	res, err := put(clientUpdateRequest(c.apiKey, endpoint, opts))
	if err != nil {
		return nil, fmt.Errorf("put: %w", err)
	}
	result := new(Client)
	if err := json.Unmarshal(res, &result); err != nil {
		if jErr, ok := err.(*json.UnmarshalTypeError); ok {
			return nil, fmt.Errorf("unmarshal field %v of type %v", jErr.Field, jErr.Type)
		}
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}
	return result, nil
}

func clientUpdateRequest(apiKey string, endpoint string, options []RequestOption) requestOptions {
	res := requestOptions{
		apiKey:   apiKey,
		endpoint: endpoint,
	}
	res.params = url.Values{}
	res.params.Add(archiveProjectsKey, strconv.FormatBool(false))
	for _, opt := range options {
		if opt.paramsProvider != nil {
			opt.paramsProvider(res.params)
		}
	}

	fields := clientUpdateFields{}
	for _, opt := range options {
		if opt.paramsProvider != nil {
			params := url.Values{}
			key := opt.paramsProvider(params)
			switch key {
			case archivedKey:
				val, _ := strconv.ParseBool(params.Get(key))
				fields.Archived = &val
			case nameKey:
				fields.Name = params.Get(key)
			}

		}
	}
	res.params.Del(archivedKey)
	res.params.Del(nameKey)
	res.fields = fields

	injectContext(&res, options)

	return res
}

// Delete existing Client.
func (c *ClientNode) Delete(workspaceID string, id string, opts ...RequestOption) (*Client,
	error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/clients/%s", c.endpoint, workspaceID, id)
	res, err := del(clientDeleteRequest(c.apiKey, endpoint, opts))
	if err != nil {
		return nil, fmt.Errorf("del: %w", err)
	}
	result := new(Client)
	if err := json.Unmarshal(res, &result); err != nil {
		if jErr, ok := err.(*json.UnmarshalTypeError); ok {
			return nil, fmt.Errorf("unmarshal field %v of type %v", jErr.Field, jErr.Type)
		}
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}
	return result, nil
}

func clientDeleteRequest(apiKey string, endpoint string, options []RequestOption) requestOptions {
	res := requestOptions{
		apiKey:   apiKey,
		endpoint: endpoint,
	}
	injectContext(&res, options)

	return res
}
