package glockify

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClientNode_All(t *testing.T) {
	ctx := context.Background()

	glock := New(clockifyTestToken)
	ws, err := glock.Workspace.All(ctx)
	require.Nil(t, err)
	require.Len(t, ws, len(workspaceNames))
	for i, workspaceName := range workspaceNames {
		require.Equal(t, workspaceName, ws[i].Name)

		cs, err := ws[i].Client.All(ctx, ClientFilter{})
		require.Nil(t, err)
		for a, clientName := range clientNames[workspaceName] {
			require.Equal(t, clientName, cs[a].Name)
		}
	}
}
