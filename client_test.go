package glockify

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type ClientTestSuite struct {
	suite.Suite
	server clientMockServer
}

type clientMockServer struct {
	baseServer *httptest.Server
	clients    []Client
}

func (c *clientMockServer) addClients(workspaceID string, count int) {
	for i := 0; i < count; i++ {
		c.clients = append(c.clients, Client{
			ID:          fmt.Sprintf("%d", i+1),
			Name:        fmt.Sprintf("Client %d", i+1),
			WorkspaceID: workspaceID,
		})
	}
}

func (s *ClientTestSuite) SetupTest() {
	testMux := mux.NewRouter()

	testMux.HandleFunc("/workspaces", workspaces)
	testMux.HandleFunc("/workspaces/{workspaceID}/clients", s.all()).Methods("GET")
	testMux.HandleFunc("/workspaces/{workspaceID}/clients/{clientID}", s.get()).Methods("GET")
	testMux.HandleFunc("/workspaces/{workspaceID}/clients", s.add()).Methods("POST")
	testMux.HandleFunc("/workspaces/{workspaceID}/clients/{clientID}", s.update()).Methods("PUT")
	testMux.HandleFunc("/workspaces/{workspaceID}/clients/{clientID}",
		s.delete()).Methods("DELETE")

	s.server = clientMockServer{
		baseServer: httptest.NewServer(testMux),
		clients:    make([]Client, 0),
	}
}

func (s *ClientTestSuite) TearDownTest() {
	s.server.baseServer.Close()
}

func (s *ClientTestSuite) all() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkAuthHeader(r) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if r.Method != "GET" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		elems := make([]string, 0)
		for _, c := range s.server.clients {
			elems = append(elems, fmt.Sprintf(`
	{
		"id": "%s",
		"name": "%s",
		"workspaceId": "%s",
		"archived": %t
  	}`, c.ID, c.Name, c.WorkspaceID, c.Archived))
		}
		_, err := fmt.Fprintf(w, "[%s]", strings.Join(elems, ","))
		if err != nil {
			log.Fatalf("%v", err)
		}
	}
}

func (s *ClientTestSuite) get() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkAuthHeader(r) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if r.Method != "GET" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		vars := mux.Vars(r)
		clientID, ok := vars["clientID"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		for _, c := range s.server.clients {
			if c.ID == clientID {
				w.WriteHeader(http.StatusOK)
				_, err := fmt.Fprintf(w, `
					{
						"id": "%s",
						"name": "%s",
						"workspaceId": "%s",
						"archived": %t
					}`, c.ID, c.Name, c.WorkspaceID, c.Archived)
				if err != nil {
					log.Fatalf("%v", err)
				}
				return
			}
		}

		w.WriteHeader(http.StatusNotFound)
	}
}

func (s *ClientTestSuite) add() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkAuthHeader(r) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if r.Method != "POST" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fields := new(ClientAddFields)
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(fields); err != nil {
			log.Fatalf("%v", err)
		}

		newClient := Client{
			ID:   fmt.Sprintf("%d", len(s.server.clients)),
			Name: fields.Name,
		}
		s.server.clients = append(s.server.clients, newClient)

		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprintf(w, `
			{
				"id": "%s",
				"name": "%s",
				"workspaceId": "%s",
				"archived": %t
			}`, newClient.ID, newClient.Name, newClient.WorkspaceID, newClient.Archived)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}
}

func (s *ClientTestSuite) update() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkAuthHeader(r) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if r.Method != "PUT" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		vars := mux.Vars(r)
		clientID, ok := vars["clientID"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fields := new(ClientUpdateFields)
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(fields); err != nil {
			log.Fatalf("%v", err)
		}

		for i, c := range s.server.clients {
			if c.ID == clientID {
				s.server.clients[i].Name = fields.Name
				w.WriteHeader(http.StatusOK)
				_, err := fmt.Fprintf(w, `
					{
						"id": "%s",
						"name": "%s",
						"workspaceId": "%s",
						"archived": %t
					}`, c.ID, fields.Name, c.WorkspaceID, c.Archived)
				if err != nil {
					log.Fatalf("%v", err)
				}
				return
			}
		}

		w.WriteHeader(http.StatusNotFound)
	}
}

func (s *ClientTestSuite) delete() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkAuthHeader(r) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if r.Method != "DELETE" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		vars := mux.Vars(r)
		clientID, ok := vars["clientID"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		foundIndex := -1
		deletedClient := Client{}
		for i, c := range s.server.clients {
			if c.ID == clientID {
				foundIndex = i
				deletedClient = c
				break
			}
		}
		if foundIndex != -1 {
			s.server.clients = append(s.server.clients[:foundIndex],
				s.server.clients[foundIndex+1:]...)
			w.WriteHeader(http.StatusOK)
			_, err := fmt.Fprintf(w, `
			{
				"id": "%s",
				"name": "%s",
				"workspaceId": "%s",
				"archived": %t
			}`, deletedClient.ID, deletedClient.Name, deletedClient.WorkspaceID,
				deletedClient.Archived)
			if err != nil {
				log.Fatalf("%v", err)
			}
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}
}

func (s *ClientTestSuite) TestAll() {
	ctx := context.Background()

	glock := New(dummyAPIKey, WithEndpoint(Endpoint{
		Base: s.server.baseServer.URL,
	}))
	ws, err := glock.Workspace.All(ctx)
	require.Nil(s.T(), err)
	require.Len(s.T(), ws, 1)

	s.server.addClients(ws[0].ID, 5)

	cs, err := ws[0].Client.All(ctx, ClientAllFilter{})
	require.Nil(s.T(), err)
	require.Len(s.T(), cs, 5)
}

func (s *ClientTestSuite) TestGet() {
	ctx := context.Background()

	glock := New(dummyAPIKey, WithEndpoint(Endpoint{
		Base: s.server.baseServer.URL,
	}))
	ws, err := glock.Workspace.All(ctx)
	require.Nil(s.T(), err)
	require.Len(s.T(), ws, 1)

	s.server.addClients(ws[0].ID, 1)

	c, err := ws[0].Client.Get(ctx, "1")
	require.Nil(s.T(), err)
	require.Equal(s.T(), "Client 1", c.Name)
}

func (s *ClientTestSuite) TestAdd() {
	ctx := context.Background()

	glock := New(dummyAPIKey, WithEndpoint(Endpoint{
		Base: s.server.baseServer.URL,
	}))
	ws, err := glock.Workspace.All(ctx)
	require.Nil(s.T(), err)
	require.Len(s.T(), ws, 1)

	wantName := "Dummy Name"

	c, err := ws[0].Client.Add(ctx, ClientAddFields{Name: wantName})
	require.Nil(s.T(), err)
	require.Equal(s.T(), wantName, c.Name)

	newC, err := ws[0].Client.Get(ctx, c.ID)
	require.Nil(s.T(), err)
	require.Equal(s.T(), wantName, newC.Name)
}

func (s *ClientTestSuite) TestUpdate() {
	ctx := context.Background()

	glock := New(dummyAPIKey, WithEndpoint(Endpoint{
		Base: s.server.baseServer.URL,
	}))
	ws, err := glock.Workspace.All(ctx)
	require.Nil(s.T(), err)
	require.Len(s.T(), ws, 1)

	s.server.addClients(ws[0].ID, 1)

	c, err := ws[0].Client.Get(ctx, "1")
	require.Nil(s.T(), err)
	require.Equal(s.T(), "Client 1", c.Name)

	uc, err := ws[0].Client.Update(ctx, c.ID, ClientUpdateFields{
		Name: "Client 2",
	}, ClientUpdateOptions{})
	require.Nil(s.T(), err)
	require.Equal(s.T(), "Client 2", uc.Name)
}

func (s *ClientTestSuite) TestDelete() {
	ctx := context.Background()

	glock := New(dummyAPIKey, WithEndpoint(Endpoint{
		Base: s.server.baseServer.URL,
	}))
	ws, err := glock.Workspace.All(ctx)
	require.Nil(s.T(), err)
	require.Len(s.T(), ws, 1)

	s.server.addClients(ws[0].ID, 5)

	c, err := ws[0].Client.Delete(ctx, "1")
	require.Nil(s.T(), err)
	require.Equal(s.T(), "Client 1", c.Name)

	cs, err := ws[0].Client.All(ctx, ClientAllFilter{})
	require.Nil(s.T(), err)
	require.Len(s.T(), cs, 4)
}

func TestClientNode(t *testing.T) {
	suite.Run(t, &ClientTestSuite{})
}
