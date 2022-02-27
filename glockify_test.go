package glockify

import "os"

var (
	clockifyTestToken = os.Getenv("CLOCKIFY_TEST_TOKEN")
	workspaceNames    = []string{"Test Workspace"}
)
