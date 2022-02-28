package glockify

import (
	"context"
	"encoding/json"
	"fmt"
)

type ClientNode struct {
	workspaceID  string
	baseEndpoint string
	apiKey       string
}

type Client struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	WorkspaceID string `json:"workspaceId,omitempty"`
	Archived    bool   `json:"archived,omitempty"`
}

type ClientFilter struct {
	Archived   bool   `schema:"archived"`
	Name       string `schema:"name"`
	Page       int    `schema:"page"`
	PageSize   int    `schema:"page-size"`
	SortColumn string `schema:"sort-column"`
	SortOrder  string `schema:"sort-order"`
}

func (c *ClientNode) All(ctx context.Context, filter ClientFilter) ([]Client, error) {
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
