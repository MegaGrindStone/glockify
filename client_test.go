package glockify

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type ClientTestSuite struct {
	suite.Suite
	server    clientMockServer
	testIndex int
}

type clientMockServer struct {
	baseServer *httptest.Server
}

func (s *ClientTestSuite) SetupTest() {
	testMux := mux.NewRouter()

	testMux.HandleFunc("/workspaces/{workspaceID}/clients", s.all()).Methods("GET")
	testMux.HandleFunc("/workspaces/{workspaceID}/clients/{clientID}", s.get()).Methods("GET")
	testMux.HandleFunc("/workspaces/{workspaceID}/clients", s.add()).Methods("POST")
	testMux.HandleFunc("/workspaces/{workspaceID}/clients/{clientID}", s.update()).Methods("PUT")
	testMux.HandleFunc("/workspaces/{workspaceID}/clients/{clientID}",
		s.delete()).Methods("DELETE")

	s.server = clientMockServer{
		baseServer: httptest.NewServer(testMux),
	}
}

func (s *ClientTestSuite) TearDownTest() {
	s.server.baseServer.Close()
}

var testsClientAll = []struct {
	name        string
	workspaceID string
	options     []RequestOption
	wantParams  url.Values
}{
	{
		name:        "Default Params",
		workspaceID: "Workspace1",
		wantParams: map[string][]string{
			"archived":   {"false"},
			"page":       {"1"},
			"page-size":  {"50"},
			"sort-order": {"DESCENDING"},
		},
	},
	{
		name:        "Set Sort Column",
		workspaceID: "Workspace1",
		options: []RequestOption{
			WithClientSortColumn(ClientSortColumnName),
		},
		wantParams: map[string][]string{
			"archived":    {"false"},
			"page":        {"1"},
			"page-size":   {"50"},
			"sort-column": {"NAME"},
			"sort-order":  {"DESCENDING"},
		},
	},
}

func (s *ClientTestSuite) all() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkAuthHeader(r) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		s.Require().Equal("GET", r.Method)

		test := testsClientAll[s.testIndex]
		path := mux.Vars(r)
		workspaceID, ok := path["workspaceID"]
		s.Require().True(ok)
		s.Require().Equal(test.workspaceID, workspaceID)
		s.Require().Equal(r.URL.Query(), test.wantParams)

		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprintf(w, `[{"id":"dummy"}]`)
		s.Require().Nil(err)
	}
}

var testsClientGet = []struct {
	name           string
	workspaceID    string
	clientID       string
	wantStatusCode int
	wantErr        bool
}{
	{
		name:           "Success",
		workspaceID:    "Workspace1",
		clientID:       "1",
		wantStatusCode: http.StatusOK,
		wantErr:        false,
	},
	{
		name:           "Not Found",
		workspaceID:    "Workspace1",
		clientID:       "2",
		wantStatusCode: http.StatusNotFound,
		wantErr:        true,
	},
}

func (s *ClientTestSuite) get() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkAuthHeader(r) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		s.Require().Equal("GET", r.Method)

		test := testsClientGet[s.testIndex]
		path := mux.Vars(r)
		workspaceID, ok := path["workspaceID"]
		s.Require().True(ok)
		s.Require().Equal(test.workspaceID, workspaceID)
		clientID, ok := path["clientID"]
		s.Require().True(ok)
		s.Require().Equal(test.clientID, clientID)

		w.WriteHeader(test.wantStatusCode)
		_, err := fmt.Fprintf(w, `{"id":"dummy"}`)
		s.Require().Nil(err)
	}
}

var testsClientAdd = []struct {
	name           string
	workspaceID    string
	clientName     string
	wantStatusCode int
	wantErr        bool
}{
	{
		name:           "Success",
		workspaceID:    "Workspace1",
		clientName:     "Client 1",
		wantStatusCode: http.StatusOK,
		wantErr:        false,
	},
	{
		name:           "Already exist",
		workspaceID:    "Workspace1",
		clientName:     "Client 2",
		wantStatusCode: http.StatusBadRequest,
		wantErr:        true,
	},
}

func (s *ClientTestSuite) add() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkAuthHeader(r) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		s.Require().Equal("POST", r.Method)
		fields := new(clientAddFields)
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(fields)
		s.Require().Nil(err)

		test := testsClientAdd[s.testIndex]
		s.Require().Equal(test.clientName, fields.Name)

		w.WriteHeader(test.wantStatusCode)
		_, err = fmt.Fprintf(w, `{"id":"dummy"}`)
		s.Require().Nil(err)
	}
}

var fl = false

var testsClientUpdate = []struct {
	name        string
	workspaceID string
	clientID    string
	options     []RequestOption
	wantParams  url.Values
	wantFields  clientUpdateFields
}{
	{
		name:        "Default",
		workspaceID: "Workspace1",
		clientID:    "1",
		wantParams: map[string][]string{
			"archive-projects": {"false"},
		},
		wantFields: clientUpdateFields{
			Archived: nil,
			Name:     "",
		},
	},
	{
		name:        "Set options",
		workspaceID: "Workspace1",
		clientID:    "2",
		options: []RequestOption{
			WithArchiveProjects(true),
			WithArchived(false),
			WithName("Dummy Name"),
		},
		wantParams: map[string][]string{
			"archive-projects": {"true"},
		},
		wantFields: clientUpdateFields{
			Archived: &fl,
			Name:     "Dummy Name",
		},
	},
}

func (s *ClientTestSuite) update() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkAuthHeader(r) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		s.Require().Equal("PUT", r.Method)

		test := testsClientUpdate[s.testIndex]

		path := mux.Vars(r)
		workspaceID, ok := path["workspaceID"]
		s.Require().True(ok)
		s.Require().Equal(test.workspaceID, workspaceID)
		clientID, ok := path["clientID"]
		s.Require().True(ok)
		s.Require().Equal(test.clientID, clientID)

		s.Require().Equal(test.wantParams, r.URL.Query())

		fields := new(clientUpdateFields)
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(fields)
		s.Require().Nil(err)
		s.Require().Equal(test.wantFields, *fields)

		w.WriteHeader(http.StatusOK)
		_, err = fmt.Fprintf(w, `{"id":"dummy"}`)
		s.Require().Nil(err)
	}
}

var testsClientDelete = []struct {
	name           string
	workspaceID    string
	clientID       string
	wantStatusCode int
	wantErr        bool
}{
	{
		name:           "Success",
		workspaceID:    "Workspace1",
		clientID:       "1",
		wantStatusCode: http.StatusOK,
		wantErr:        false,
	},
	{
		name:           "Not Found",
		workspaceID:    "Workspace1",
		clientID:       "2",
		wantStatusCode: http.StatusNotFound,
		wantErr:        true,
	},
}

func (s *ClientTestSuite) delete() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkAuthHeader(r) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		s.Require().Equal("DELETE", r.Method)

		test := testsClientDelete[s.testIndex]

		path := mux.Vars(r)
		workspaceID, ok := path["workspaceID"]
		s.Require().True(ok)
		s.Require().Equal(test.workspaceID, workspaceID)
		clientID, ok := path["clientID"]
		s.Require().True(ok)
		s.Require().Equal(test.clientID, clientID)

		w.WriteHeader(test.wantStatusCode)
		_, err := fmt.Fprintf(w, `{"id":"dummy"}`)
		s.Require().Nil(err)
	}
}

func (s *ClientTestSuite) TestAll() {
	glock := New(dummyAPIKey, WithEndpoint(Endpoint{
		Base: s.server.baseServer.URL,
	}))

	for index, tc := range testsClientAll {
		s.Run(tc.name, func() {
			s.testIndex = index
			_, err := glock.Client.All(tc.workspaceID, tc.options...)
			s.Require().Nil(err)
		})
	}
}

func (s *ClientTestSuite) TestGet() {
	glock := New(dummyAPIKey, WithEndpoint(Endpoint{
		Base: s.server.baseServer.URL,
	}))

	for index, tc := range testsClientGet {
		s.Run(tc.name, func() {
			s.testIndex = index
			_, err := glock.Client.Get(tc.workspaceID, tc.clientID)
			if tc.wantErr {
				s.Require().NotNil(err)
			} else {
				s.Require().Nil(err)
			}
		})
	}
}

func (s *ClientTestSuite) TestAdd() {
	glock := New(dummyAPIKey, WithEndpoint(Endpoint{
		Base: s.server.baseServer.URL,
	}))

	for index, tc := range testsClientAdd {
		s.Run(tc.name, func() {
			s.testIndex = index
			_, err := glock.Client.Add(tc.workspaceID, tc.clientName)
			if tc.wantErr {
				s.Require().NotNil(err)
			} else {
				s.Require().Nil(err)
			}
		})
	}
}

func (s *ClientTestSuite) TestUpdate() {
	glock := New(dummyAPIKey, WithEndpoint(Endpoint{
		Base: s.server.baseServer.URL,
	}))

	for index, tc := range testsClientUpdate {
		s.Run(tc.name, func() {
			s.testIndex = index
			_, err := glock.Client.Update(tc.workspaceID, tc.clientID, tc.options...)
			s.Require().Nil(err)
		})
	}
}

func (s *ClientTestSuite) TestDelete() {
	glock := New(dummyAPIKey, WithEndpoint(Endpoint{
		Base: s.server.baseServer.URL,
	}))

	for index, tc := range testsClientDelete {
		s.Run(tc.name, func() {
			s.testIndex = index
			_, err := glock.Client.Delete(tc.workspaceID, tc.clientID)
			if tc.wantErr {
				s.Require().NotNil(err)
			} else {
				s.Require().Nil(err)
			}
		})
	}
}

func TestClientNode(t *testing.T) {
	suite.Run(t, &ClientTestSuite{})
}
