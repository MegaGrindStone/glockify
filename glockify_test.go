package glockify

import "os"

var (
	clockifyTestToken = os.Getenv("CLOCKIFY_TEST_TOKEN")
	workspaceNames    = []string{"Test Workspace"}
	clientNames       = map[string][]string{
		"Test Workspace": {
			"Test Client 1",
			"Test Client 2",
		},
	}
)
