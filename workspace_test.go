package glockify

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGlockify_Workspaces(t *testing.T) {
	glock := New(clockifyTestToken)
	ws, err := glock.Workspace.All(context.Background())
	require.Nil(t, err)
	require.Len(t, ws, len(workspaceNames))
	for i, workspaceName := range workspaceNames {
		require.Equal(t, workspaceName, ws[i].Name)
	}
}
