package glockify

import (
	"context"
	"encoding/json"
	"fmt"
)

// ClientNode manipulating Client resource.
type ClientNode struct {
	endpoint string
	apiKey   string
}

// Client wraps Clockify's client resource.
// See: https://clockify.me/developers-api#tag-Client
type Client struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	WorkspaceID string `json:"workspaceId,omitempty"`
	Archived    bool   `json:"archived,omitempty"`
}

// ClientAllFilter is used for All request.
type ClientAllFilter struct {
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

	// SortColumn possible values is "NAME"
	SortColumn string `schema:"sort-column"`

	// SortOrder possible values is "ASCENDING", "DESCENDING"
	SortOrder string `schema:"sort-order"`
}

// ClientAddFields is used for Add request.
// See: https://clockify.me/developers-api#tag-Client
type ClientAddFields struct {
	Name string `json:"name"`
}

// ClientUpdateFields is used for Update request.
// See: https://clockify.me/developers-api#tag-Client
type ClientUpdateFields struct {
	Archived bool   `json:"archived"`
	Name     string `json:"name"`
}

// ClientUpdateOptions is used for Update request.
type ClientUpdateOptions struct {
	// ArchiveProjects controls whether archiving client will result in archiving
	// all projects of given client.
	ArchiveProjects bool `scheme:"archive-projects"`
}

// All get all Client resource based on filter given.
func (c *ClientNode) All(ctx context.Context, workspaceID string,
	filter ClientAllFilter) ([]Client, error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/clients", c.endpoint, workspaceID)
	res, err := get(ctx, c.apiKey, filter, endpoint)
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

// Get one Client by its id.
func (c *ClientNode) Get(ctx context.Context, workspaceID string, id string) (*Client, error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/clients/%s", c.endpoint, workspaceID, id)
	res, err := get(ctx, c.apiKey, nil, endpoint)
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

// Add create new Client based on fields given.
func (c *ClientNode) Add(ctx context.Context, workspaceID string, fields ClientAddFields) (*Client,
	error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/clients", c.endpoint, workspaceID)
	res, err := post(ctx, c.apiKey, nil, fields, endpoint)
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

// Update existing Client based on fields and options given.
func (c *ClientNode) Update(ctx context.Context, workspaceID string, id string,
	fields ClientUpdateFields, options ClientUpdateOptions) (*Client, error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/clients/%s", c.endpoint, workspaceID, id)
	res, err := put(ctx, c.apiKey, options, fields, endpoint)
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

// Delete existing Client.
func (c *ClientNode) Delete(ctx context.Context, workspaceID string, id string) (*Client, error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/clients/%s", c.endpoint, workspaceID, id)
	res, err := del(ctx, c.apiKey, endpoint)
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
