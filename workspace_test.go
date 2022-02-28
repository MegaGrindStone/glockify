package glockify

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http/httptest"
	"testing"
)

type WorkspaceTestSuite struct {
	suite.Suite
	server workspaceMockServer
}

type workspaceMockServer struct {
	baseServer *httptest.Server
}

func (s *WorkspaceTestSuite) SetupTest() {
	testMux := mux.NewRouter()

	testMux.HandleFunc("/workspaces", workspaces)

	s.server = workspaceMockServer{
		baseServer: httptest.NewServer(testMux),
	}
}

func (s *WorkspaceTestSuite) TearDownTest() {
	s.server.baseServer.Close()
}

func (s *WorkspaceTestSuite) TestAll() {
	ctx := context.Background()

	glock := New(dummyAPIKey, WithEndpoint(Endpoint{
		Base: s.server.baseServer.URL,
	}))

	ws, err := glock.Workspace.All(ctx)
	require.Nil(s.T(), err)
	require.Len(s.T(), ws, 1)
}

func TestWorkspaceNode(t *testing.T) {
	suite.Run(t, &WorkspaceTestSuite{})
}
