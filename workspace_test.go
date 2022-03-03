package glockify

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type WorkspaceTestSuite struct {
	suite.Suite
	server    workspaceMockServer
	testIndex int
}

type workspaceMockServer struct {
	baseServer *httptest.Server
}

func (s *WorkspaceTestSuite) SetupTest() {
	testMux := mux.NewRouter()

	testMux.HandleFunc("/workspaces", s.all())

	s.server = workspaceMockServer{
		baseServer: httptest.NewServer(testMux),
	}
}

func (s *WorkspaceTestSuite) TearDownTest() {
	s.server.baseServer.Close()
}

var testsWorkspaceAll = []struct {
	name string
}{
	{
		name: "Default Params",
	},
}

func (s *WorkspaceTestSuite) all() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkAuthHeader(r) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		s.Require().Equal("GET", r.Method)

		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprintf(w, `[{"id":"dummy"}]`)
		s.Require().Nil(err)
	}
}

func (s *WorkspaceTestSuite) TestAll() {
	glock := New(dummyAPIKey, WithEndpoint(Endpoint{
		Base: s.server.baseServer.URL,
	}))

	for index, tc := range testsWorkspaceAll {
		s.Run(tc.name, func() {
			s.testIndex = index
			_, err := glock.Workspace.All()
			s.Require().Nil(err)
		})
	}
}

func TestWorkspaceNode(t *testing.T) {
	suite.Run(t, &WorkspaceTestSuite{})
}
