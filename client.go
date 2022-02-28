package glockify

import (
	"context"
	"encoding/json"
	"fmt"
)

// ClientNode incorporating every Client resource request.
type ClientNode struct {
	workspaceID  string
	baseEndpoint string
	apiKey       string
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
// See: https://clockify.me/developers-api#tag-Client
type ClientAllFilter struct {
	Archived   bool   `schema:"archived"`
	Name       string `schema:"name"`
	Page       int    `schema:"page"`
	PageSize   int    `schema:"page-size"`
	SortColumn string `schema:"sort-column"`
	SortOrder  string `schema:"sort-order"`
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
// See: https://clockify.me/developers-api#tag-Client
type ClientUpdateOptions struct {
	ArchiveProjects bool `scheme:"archive-projects"`
}

// All get all Client resource based on filter given.
func (c *ClientNode) All(ctx context.Context, filter ClientAllFilter) ([]Client, error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/clients", c.baseEndpoint, c.workspaceID)
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

// Get one client by its id.
func (c *ClientNode) Get(ctx context.Context, id string) (*Client, error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/clients/%s", c.baseEndpoint, c.workspaceID, id)
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
func (c *ClientNode) Add(ctx context.Context, fields ClientAddFields) (*Client, error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/clients", c.baseEndpoint, c.workspaceID)
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
func (c *ClientNode) Update(ctx context.Context, id string, fields ClientUpdateFields,
	options ClientUpdateOptions) (*Client, error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/clients/%s", c.baseEndpoint, c.workspaceID, id)
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
func (c *ClientNode) Delete(ctx context.Context, id string) (*Client, error) {
	endpoint := fmt.Sprintf("%s/workspaces/%s/clients/%s", c.baseEndpoint, c.workspaceID, id)
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
