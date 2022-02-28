package glockify

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/schema"
	"io/ioutil"
	"log"
	"net/http"
)

type Glockify struct {
	Workspace WorkspaceNode
}

type Endpoint struct {
	Base    string
	TimeOff string
	Report  string
}

type Option func(*Glockify)

const (
	defaultBaseEndpoint    = "https://api.clockify.me/api/v1"
	defaultTimeOffEndpoint = "https://reports.api.clockify.me/v1"
	defaultReportEndpoint  = "https://pto.api.clockify.me/v1"
)

func New(apiKey string, opts ...Option) *Glockify {
	g := &Glockify{
		Workspace: WorkspaceNode{
			endpoint: Endpoint{
				Base:    defaultBaseEndpoint,
				TimeOff: defaultTimeOffEndpoint,
				Report:  defaultReportEndpoint,
			},
			apiKey: apiKey,
		},
	}

	for _, opt := range opts {
		opt(g)
	}

	return g
}

func WithEndpoint(endpoint Endpoint) Option {
	return func(g *Glockify) {
		if endpoint.Base != "" {
			g.Workspace.endpoint.Base = endpoint.Base
		}
		if endpoint.TimeOff != "" {
			g.Workspace.endpoint.TimeOff = endpoint.TimeOff
		}
		if endpoint.Report != "" {
			g.Workspace.endpoint.Report = endpoint.Report
		}
	}
}

func get(ctx context.Context, apiKey string, params interface{}, endpoint string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("X-Api-Key", apiKey)
	if params != nil {
		encoder := schema.NewEncoder()
		if err := encoder.Encode(params, req.URL.Query()); err != nil {
			return nil, fmt.Errorf("scheme encode: %w", err)
		}
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf(fmt.Errorf("close body: %w", err).Error())
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("http error: status code %d", resp.StatusCode))
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	return respBytes, nil
}
