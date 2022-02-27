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
	require.NotEmpty(t, ws)
}
