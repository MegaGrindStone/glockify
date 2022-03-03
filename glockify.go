package glockify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// Glockify is an entry point to access Clockify API.
type Glockify struct {
	Workspace WorkspaceNode
	Client    ClientNode
	Project   ProjectNode
	Task      TaskNode

	apiKey string
}

// Endpoint specify main endpoints in Clockify.
type Endpoint struct {
	Base    string
	TimeOff string
	Report  string
}

// Option control parameter that can given when creating new Glockify.
type Option func(*Glockify)

const (
	defaultBaseEndpoint    = "https://api.clockify.me/api/v1"
	defaultTimeOffEndpoint = "https://reports.api.clockify.me/v1"
	defaultReportEndpoint  = "https://pto.api.clockify.me/v1"
)

// New instantiate Glockify with apiKey given.
func New(apiKey string, opts ...Option) *Glockify {
	g := &Glockify{
		apiKey: apiKey,
	}
	g.setupNode(Endpoint{
		Base:    defaultBaseEndpoint,
		TimeOff: defaultTimeOffEndpoint,
		Report:  defaultReportEndpoint,
	})

	for _, opt := range opts {
		opt(g)
	}

	return g
}

func (g *Glockify) setupNode(endpoint Endpoint) {
	g.Workspace = WorkspaceNode{
		endpoint: endpoint.Base,
		apiKey:   g.apiKey,
	}
	g.Client = ClientNode{
		endpoint: endpoint.Base,
		apiKey:   g.apiKey,
	}
	g.Project = ProjectNode{
		endpoint: endpoint.Base,
		apiKey:   g.apiKey,
	}
}

// WithEndpoint set endpoint when creating new Glockify.
func WithEndpoint(endpoint Endpoint) Option {
	return func(g *Glockify) {
		defaultEndpoint := Endpoint{
			Base:    defaultBaseEndpoint,
			TimeOff: defaultTimeOffEndpoint,
			Report:  defaultReportEndpoint,
		}
		if endpoint.Base != "" {
			defaultEndpoint.Base = endpoint.Base
		}
		if endpoint.TimeOff != "" {
			defaultEndpoint.TimeOff = endpoint.TimeOff
		}
		if endpoint.Report != "" {
			defaultEndpoint.Report = endpoint.Report
		}
		g.setupNode(defaultEndpoint)
	}
}

type requestOptions struct {
	ctx      context.Context
	apiKey   string
	endpoint string
	params   url.Values
	fields   interface{}
}

func get(opt requestOptions) ([]byte, error) {
	req, err := http.NewRequestWithContext(opt.ctx, "GET",
		opt.endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("X-Api-Key", opt.apiKey)
	if opt.params != nil {
		req.URL.RawQuery = opt.params.Encode()
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
		return nil, fmt.Errorf("http error: status code %d", resp.StatusCode)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	return respBytes, nil
}

func post(opt requestOptions) ([]byte, error) {
	buf := new(bytes.Buffer)
	if opt.fields != nil {
		bodyJSON, err := json.Marshal(opt.fields)
		if err != nil {
			return nil, fmt.Errorf("json marshal: %w", err)
		}
		buf = bytes.NewBuffer(bodyJSON)
	}
	req, err := http.NewRequestWithContext(opt.ctx, "POST", opt.endpoint, buf)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("X-Api-Key", opt.apiKey)

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

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("http error: status code %d", resp.StatusCode)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	return respBytes, nil
}

func put(opt requestOptions) ([]byte, error) {
	buf := new(bytes.Buffer)
	if opt.fields != nil {
		bodyJSON, err := json.Marshal(opt.fields)
		if err != nil {
			return nil, fmt.Errorf("json marshal: %w", err)
		}
		buf = bytes.NewBuffer(bodyJSON)
	}
	req, err := http.NewRequestWithContext(opt.ctx, "PUT", opt.endpoint, buf)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("X-Api-Key", opt.apiKey)
	if opt.params != nil {
		req.URL.RawQuery = opt.params.Encode()
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
		return nil, fmt.Errorf("http error: status code %d", resp.StatusCode)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	return respBytes, nil
}

func patch(opt requestOptions) ([]byte, error) {
	buf := new(bytes.Buffer)
	if opt.fields != nil {
		bodyJSON, err := json.Marshal(opt.fields)
		if err != nil {
			return nil, fmt.Errorf("json marshal: %w", err)
		}
		buf = bytes.NewBuffer(bodyJSON)
	}
	req, err := http.NewRequestWithContext(opt.ctx, "PATCH", opt.endpoint, buf)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("X-Api-Key", opt.apiKey)
	if opt.params != nil {
		req.URL.RawQuery = opt.params.Encode()
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
		return nil, fmt.Errorf("http error: status code %d", resp.StatusCode)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	return respBytes, nil
}

func del(opt requestOptions) ([]byte, error) {
	req, err := http.NewRequestWithContext(opt.ctx, "DELETE", opt.endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("X-Api-Key", opt.apiKey)
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
		return nil, fmt.Errorf("http error: status code %d", resp.StatusCode)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	return respBytes, nil
}

type contextOptions struct {
	ctx context.Context
}

func injectContext(requestOptions *requestOptions, opts []RequestOption) {
	co := &contextOptions{ctx: context.Background()}
	for _, opt := range opts {
		if opt.contextProvider != nil {
			opt.contextProvider(co)
		}
	}
	requestOptions.ctx = co.ctx
}

// WithContext set request context. Default to context.Background.
func WithContext(ctx context.Context) RequestOption {
	return RequestOption{
		contextProvider: func(o *contextOptions) {
			o.ctx = ctx
		},
	}
}

const (
	defaultPage      = 1
	defaultPageSize  = 50
	defaultSortOrder = SortOrderDescending

	arraySeparator = "|"

	pageKey       = "page"
	pageSizeKey   = "page-size"
	sortColumnKey = "sort-column"
	sortOrderKey  = "sort-order"
	archivedKey   = "archived"
	nameKey       = "name"
	billableKey   = "billable"
)

type SortOrderValue string

// Possible value for SortOrderValue
const (
	SortOrderAscending  SortOrderValue = "ASCENDING"
	SortOrderDescending                = "DESCENDING"
)

// WithPage set request's page. Default to 1.
func WithPage(page int) RequestOption {
	return RequestOption{
		paramsProvider: func(v url.Values) string {
			v.Set(pageKey, strconv.Itoa(page))
			return pageKey
		},
	}
}

// WithPageSize set length of items returned from request. Default to 50.
// Maximum value is 5000.
func WithPageSize(pageSize int) RequestOption {
	return RequestOption{
		paramsProvider: func(v url.Values) string {
			v.Set(pageSizeKey, strconv.Itoa(pageSize))
			return pageSizeKey
		},
	}
}

// WithSortOrder set sorting behaviour. Default to SortOrderDescending.
func WithSortOrder(sortOrder SortOrderValue) RequestOption {
	return RequestOption{
		paramsProvider: func(v url.Values) string {
			v.Set(sortOrderKey, string(sortOrder))
			return sortOrderKey
		},
	}
}

// WithArchived when applied to All request, its filter response by active/archived state,
// and defaulted to false. When applied to Update request, its set entity state
// to active/archived, and defaulted to not set.
func WithArchived(archived bool) RequestOption {
	return RequestOption{
		paramsProvider: func(v url.Values) string {
			v.Set(archivedKey, strconv.FormatBool(archived))
			return archivedKey
		},
	}
}

// WithName when applied to All request, its filter response by its name.
// When applied to Update its set entity name.
func WithName(name string) RequestOption {
	return RequestOption{
		paramsProvider: func(v url.Values) string {
			v.Set(nameKey, name)
			return nameKey
		},
	}
}

// WithBillable when applied to All, its filter response by its billable state.
// When applied to Update its set entity billable state.
func WithBillable(billable bool) RequestOption {
	return RequestOption{
		paramsProvider: func(v url.Values) string {
			v.Set(billableKey, strconv.FormatBool(billable))
			return billableKey
		},
	}
}

type RequestOption struct {
	contextProvider func(*contextOptions)
	paramsProvider  func(url.Values) string
}
